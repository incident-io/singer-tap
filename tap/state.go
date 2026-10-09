package tap

import (
	"time"

	"github.com/pkg/errors"
)

// BookmarkKey is the replication key, and the key under which each incremental stream's
// bookmark is stored in the state.
const BookmarkKey = "updated_at"

// State is the Singer state passed in with --state, and emitted in STATE messages.
//
// It has the shape {"bookmarks": {"follow_ups": {"updated_at": "2026-10-08T12:00:00Z"}}}.
type State struct {
	Bookmarks map[string]map[string]string `json:"bookmarks"`
}

// GetBookmark returns the bookmark for the stream, or nil if there isn't one.
func (s *State) GetBookmark(stream string) (*time.Time, error) {
	value, ok := s.Bookmarks[stream][BookmarkKey]
	if !ok {
		return nil, nil
	}

	bookmark, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, errors.Wrapf(err, "parsing %s bookmark for stream %s", BookmarkKey, stream)
	}

	return &bookmark, nil
}

// SetBookmark sets the bookmark for the stream.
func (s *State) SetBookmark(stream string, bookmark time.Time) {
	if s.Bookmarks == nil {
		s.Bookmarks = map[string]map[string]string{}
	}

	s.Bookmarks[stream] = map[string]string{
		BookmarkKey: bookmark.UTC().Format(time.RFC3339),
	}
}

// UpdatedSinceFilter builds the updated_at filter for a list request from a bookmark.
//
// The API documents its timestamp filters with dates only (e.g. "2025-01-01") and doesn't
// say which timezone it reads them in. We send the date of the day before the bookmark, so
// we fetch everything updated since the bookmark whatever the timezone, at the cost of
// re-fetching up to two days of changes. Targets upsert on id, so the overlap is harmless.
func UpdatedSinceFilter(since time.Time) *map[string][]string {
	return &map[string][]string{
		"gte": {since.UTC().Add(-24 * time.Hour).Format(time.DateOnly)},
	}
}
