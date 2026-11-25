package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"time"

	"github.com/incident-io/singer-tap/client"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("New", func() {
	var (
		ctx        context.Context
		server     *httptest.Server
		requests   atomic.Int32
		statuses   []int
		retryAfter string
		mu         sync.Mutex
		arrivals   []time.Time
	)

	BeforeEach(func() {
		ctx = context.Background()
		requests.Store(0)
		statuses = nil
		retryAfter = ""
		arrivals = nil

		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mu.Lock()
			arrivals = append(arrivals, time.Now())
			mu.Unlock()

			attempt := int(requests.Add(1)) - 1
			status := statuses[min(attempt, len(statuses)-1)]

			w.Header().Set("Content-Type", "application/json")
			if status == http.StatusTooManyRequests && retryAfter != "" {
				w.Header().Set("Retry-After", retryAfter)
			}
			w.WriteHeader(status)
			if status == http.StatusOK {
				_, _ = w.Write([]byte(`{"severities":[]}`))
			} else {
				_, _ = w.Write([]byte("upstream connect error"))
			}
		}))
		DeferCleanup(server.Close)
	})

	listSeverities := func() error {
		cl, err := client.New(ctx, "an-api-key", server.URL, "test")
		Expect(err).NotTo(HaveOccurred())

		_, err = cl.SeveritiesV1ListWithResponse(ctx)
		return err
	}

	When("the first response succeeds", func() {
		BeforeEach(func() {
			statuses = []int{http.StatusOK}
		})

		It("makes one request", func() {
			Expect(listSeverities()).To(Succeed())
			Expect(requests.Load()).To(BeNumerically("==", 1))
		})
	})

	When("the API returns a 502 and then recovers", func() {
		BeforeEach(func() {
			statuses = []int{http.StatusBadGateway, http.StatusBadGateway, http.StatusOK}
		})

		It("retries until it succeeds", func() {
			Expect(listSeverities()).To(Succeed())
			Expect(requests.Load()).To(BeNumerically("==", 3))
		})
	})

	When("the API keeps returning a 502", func() {
		BeforeEach(func() {
			statuses = []int{http.StatusBadGateway}
		})

		It("gives up after the retries and reports the status and body", func() {
			err := listSeverities()
			Expect(err).To(MatchError(ContainSubstring("status 502: upstream connect error")))
			Expect(requests.Load()).To(BeNumerically("==", 4))
		})
	})

	When("the API returns a 429", func() {
		BeforeEach(func() {
			statuses = []int{http.StatusTooManyRequests, http.StatusOK}
		})

		It("waits before it retries", func() {
			Expect(listSeverities()).To(Succeed())
			Expect(arrivals).To(HaveLen(2))
			Expect(arrivals[1].Sub(arrivals[0])).To(BeNumerically(">=", time.Second))
		})

		When("the response has a Retry-After header", func() {
			BeforeEach(func() {
				retryAfter = "2"
			})

			It("waits for that long before it retries", func() {
				Expect(listSeverities()).To(Succeed())
				Expect(arrivals).To(HaveLen(2))
				Expect(arrivals[1].Sub(arrivals[0])).To(BeNumerically(">=", 2*time.Second))
			})
		})
	})

	When("the API returns a 400", func() {
		BeforeEach(func() {
			statuses = []int{http.StatusBadRequest}
		})

		It("fails without retrying", func() {
			err := listSeverities()
			Expect(err).To(MatchError(ContainSubstring("status 400")))
			Expect(requests.Load()).To(BeNumerically("==", 1))
		})
	})
})
