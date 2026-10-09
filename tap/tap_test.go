package tap_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"time"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
	"github.com/incident-io/singer-tap/tap"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sync", func() {
	var (
		ctx    context.Context
		server *fakeAPI
		out    *bytes.Buffer
		state  *tap.State
	)

	// sync runs the tap against the fake API, for the follow_ups and actions streams only.
	sync := func(catalog *tap.Catalog) error {
		cl, err := client.New(ctx, "test-api-key", server.URL, "test")
		Expect(err).NotTo(HaveOccurred())

		return tap.Sync(ctx, kitlog.NewNopLogger(), tap.NewOutputLogger(out), cl, catalog, state)
	}

	BeforeEach(func() {
		ctx = context.Background()
		out = &bytes.Buffer{}
		state = nil
		server = newFakeAPI()
		DeferCleanup(server.Close)
	})

	Context("without state", func() {
		It("fetches every page with no updated_at filter", func() {
			Expect(sync(incrementalCatalog())).To(Succeed())

			for _, path := range []string{"/v3/follow_ups", "/v3/actions"} {
				queries := server.Requests(path)
				Expect(queries).To(HaveLen(2), path)
				Expect(queries[0].Has("updated_at[gte]")).To(BeFalse(), path)
				Expect(queries[1].Has("updated_at[gte]")).To(BeFalse(), path)
				Expect(queries[1].Get("after")).To(Equal("cursor-1"), path)
			}

			Expect(recordIDs(out, "follow_ups")).To(ConsistOf("follow-up-1", "follow-up-2"))
			Expect(recordIDs(out, "actions")).To(ConsistOf("action-1", "action-2"))
		})

		It("emits a STATE after each stream, bookmarked at the time the stream started", func() {
			before := time.Now().UTC().Truncate(time.Second)
			Expect(sync(incrementalCatalog())).To(Succeed())
			after := time.Now().UTC()

			states := stateMessages(out)
			Expect(states).To(HaveLen(2))

			last := states[len(states)-1]
			for _, stream := range []string{"follow_ups", "actions"} {
				bookmark, err := last.GetBookmark(stream)
				Expect(err).NotTo(HaveOccurred())
				Expect(bookmark).NotTo(BeNil(), stream)
				Expect(*bookmark).To(BeTemporally(">=", before), stream)
				Expect(*bookmark).To(BeTemporally("<=", after), stream)
			}
		})

		It("emits each STATE after that stream's records", func() {
			Expect(sync(incrementalCatalog())).To(Succeed())

			var types []string
			for _, msg := range messages(out) {
				switch msg["type"] {
				case "RECORD":
					types = append(types, fmt.Sprintf("RECORD:%s", msg["stream"]))
				case "STATE":
					types = append(types, "STATE")
				}
			}

			Expect(types).To(Equal([]string{
				"RECORD:actions", "RECORD:actions", "STATE",
				"RECORD:follow_ups", "RECORD:follow_ups", "STATE",
			}))
		})
	})

	Context("with state", func() {
		BeforeEach(func() {
			state = &tap.State{
				Bookmarks: map[string]map[string]string{
					"follow_ups":   {"updated_at": "2026-10-08T12:00:00Z"},
					"other_stream": {"updated_at": "2026-01-01T00:00:00Z"},
				},
			}
		})

		It("filters every follow-ups page on updated_at, from the day before the bookmark", func() {
			Expect(sync(incrementalCatalog())).To(Succeed())

			queries := server.Requests("/v3/follow_ups")
			Expect(queries).To(HaveLen(2))
			for _, query := range queries {
				Expect(query["updated_at[gte]"]).To(Equal([]string{"2026-10-07"}))
			}
		})

		It("does not filter streams with no bookmark", func() {
			Expect(sync(incrementalCatalog())).To(Succeed())

			for _, query := range server.Requests("/v3/actions") {
				Expect(query.Has("updated_at[gte]")).To(BeFalse())
			}
		})

		It("advances the bookmark and keeps bookmarks for other streams", func() {
			before := time.Now().UTC().Truncate(time.Second)
			Expect(sync(incrementalCatalog())).To(Succeed())

			states := stateMessages(out)
			last := states[len(states)-1]

			bookmark, err := last.GetBookmark("follow_ups")
			Expect(err).NotTo(HaveOccurred())
			Expect(*bookmark).To(BeTemporally(">=", before))

			Expect(last.Bookmarks["other_stream"]).To(Equal(map[string]string{"updated_at": "2026-01-01T00:00:00Z"}))
		})

		It("errors on a bookmark that isn't a timestamp", func() {
			state.Bookmarks["follow_ups"]["updated_at"] = "yesterday"

			Expect(sync(incrementalCatalog())).To(MatchError(ContainSubstring("parsing updated_at bookmark for stream follow_ups")))
		})
	})

	Context("when a stream fails partway through", func() {
		BeforeEach(func() {
			server.FailSecondPage("/v3/follow_ups")
		})

		It("does not emit a bookmark for that stream", func() {
			Expect(sync(incrementalCatalog())).NotTo(Succeed())

			for _, emitted := range stateMessages(out) {
				Expect(emitted.Bookmarks).NotTo(HaveKey("follow_ups"))
			}
		})
	})

	Context("when the catalog deselects the replication key", func() {
		It("still emits updated_at and id", func() {
			catalog := incrementalCatalog()
			for _, entry := range catalog.Streams {
				for idx, metadata := range *entry.Metadata {
					if len(metadata.Breadcrumb) == 2 && (metadata.Breadcrumb[1] == "updated_at" || metadata.Breadcrumb[1] == "id") {
						(*entry.Metadata)[idx].Metadata.Selected = new(bool)
					}
				}
			}

			Expect(sync(catalog)).To(Succeed())

			for _, msg := range messages(out) {
				if msg["type"] == "RECORD" {
					Expect(msg["record"]).To(HaveKey("updated_at"))
					Expect(msg["record"]).To(HaveKey("id"))
				}
			}
		})
	})
})

var _ = Describe("Discover", func() {
	It("marks follow_ups and actions as incremental on updated_at", func() {
		catalog := discover()

		for _, entry := range catalog.Streams {
			top := (*entry.Metadata)[0].Metadata
			switch entry.Stream {
			case "follow_ups", "actions":
				Expect(top.ForcedReplicationMethod).To(Equal("INCREMENTAL"), entry.Stream)
				Expect(top.ValidReplicationKeys).To(Equal([]string{"updated_at"}), entry.Stream)
			default:
				Expect(top.ForcedReplicationMethod).To(Equal("FULL_TABLE"), entry.Stream)
			}
		}
	})
})

var _ = Describe("UpdatedSinceFilter", func() {
	It("uses the UTC date of the day before", func() {
		// 2026-10-08T01:00 in UTC+3 is 2026-10-07T22:00 UTC, so the day before is the 6th.
		since := time.Date(2026, 10, 8, 1, 0, 0, 0, time.FixedZone("UTC+3", 3*60*60))

		Expect(*tap.UpdatedSinceFilter(since)).To(Equal(map[string][]string{"gte": {"2026-10-06"}}))
	})
})

// discover returns the tap's default catalog.
func discover() *tap.Catalog {
	out := &bytes.Buffer{}
	Expect(tap.Discover(context.Background(), kitlog.NewNopLogger(), tap.NewOutputLogger(out))).To(Succeed())

	var catalog tap.Catalog
	Expect(json.Unmarshal(out.Bytes(), &catalog)).To(Succeed())

	return &catalog
}

// incrementalCatalog returns the default catalog, limited to the incremental streams.
func incrementalCatalog() *tap.Catalog {
	catalog := discover()

	var entries []tap.CatalogEntry
	for _, entry := range catalog.Streams {
		if entry.Stream == "follow_ups" || entry.Stream == "actions" {
			entries = append(entries, entry)
		}
	}
	catalog.Streams = entries

	return catalog
}

func messages(out *bytes.Buffer) []map[string]any {
	var msgs []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if line == "" {
			continue
		}

		var msg map[string]any
		Expect(json.Unmarshal([]byte(line), &msg)).To(Succeed())
		msgs = append(msgs, msg)
	}

	return msgs
}

func recordIDs(out *bytes.Buffer, stream string) []string {
	var ids []string
	for _, msg := range messages(out) {
		if msg["type"] == "RECORD" && msg["stream"] == stream {
			ids = append(ids, msg["record"].(map[string]any)["id"].(string))
		}
	}

	return ids
}

func stateMessages(out *bytes.Buffer) []tap.State {
	var states []tap.State
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		var msg struct {
			Type  string    `json:"type"`
			Value tap.State `json:"value"`
		}
		Expect(json.Unmarshal([]byte(line), &msg)).To(Succeed())

		if msg.Type == "STATE" {
			states = append(states, msg.Value)
		}
	}

	return states
}

// fakeAPI serves two pages of follow-ups and actions, and records the query of every
// request it receives.
type fakeAPI struct {
	*httptest.Server

	mu         sync.Mutex
	requests   map[string][]url.Values
	failSecond map[string]bool
}

func newFakeAPI() *fakeAPI {
	f := &fakeAPI{
		requests:   map[string][]url.Values{},
		failSecond: map[string]bool{},
	}
	f.Server = httptest.NewServer(http.HandlerFunc(f.serve))

	return f
}

func (f *fakeAPI) Requests(path string) []url.Values {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.requests[path]
}

// FailSecondPage makes the second page of the path return a 400, which the client doesn't
// retry.
func (f *fakeAPI) FailSecondPage(path string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.failSecond[path] = true
}

func (f *fakeAPI) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	f.requests[r.URL.Path] = append(f.requests[r.URL.Path], r.URL.Query())
	failSecond := f.failSecond[r.URL.Path]
	f.mu.Unlock()

	var key, prefix string
	switch r.URL.Path {
	case "/v3/follow_ups":
		key, prefix = "follow_ups", "follow-up"
	case "/v3/actions":
		key, prefix = "actions", "action"
	default:
		http.NotFound(w, r)
		return
	}

	// The first page returns a cursor, the second doesn't.
	page, after := 1, any("cursor-1")
	if r.URL.Query().Get("after") == "cursor-1" {
		page, after = 2, nil
	}

	if page == 2 && failSecond {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"type": "validation_error"}`)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	Expect(json.NewEncoder(w).Encode(map[string]any{
		key: []map[string]any{{
			"id":          fmt.Sprintf("%s-%d", prefix, page),
			"incident_id": "incident-1",
			"status":      "outstanding",
			"title":       "A title",
			"description": "A description",
			"created_at":  "2026-10-01T00:00:00Z",
			"updated_at":  "2026-10-08T00:00:00Z",
		}},
		"pagination_meta": map[string]any{
			"after":     after,
			"page_size": 250,
		},
	})).To(Succeed())
}
