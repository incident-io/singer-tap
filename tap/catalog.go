package tap

import "github.com/incident-io/singer-tap/model"

// A catalog can contain several streams or "entries"
type CatalogEntry struct {
	// Name of the stream
	Stream string `json:"stream"`

	// Unique identifier for this stream
	// Allows for multiple sources that have duplicate stream names
	TapStreamID string `json:"tap_stream_id"`

	// The given schema for this stream
	Schema model.Schema `json:"schema"`

	// Optional metadata for this stream
	// Metadata *[]Metadata `json:"metadata,omitempty"`
}

// Actual catalog that we export
// contains an array of all our streams
type Catalog struct {
	Streams []CatalogEntry `json:"streams"`
}

func NewCatalog(streams map[string]Stream) *Catalog {
	entries := []CatalogEntry{}

	for name, stream := range streams {
		catalogEntry := CatalogEntry{
			Stream:      name,
			TapStreamID: name,
			Schema:      *stream.Output().Schema,
		}

		entries = append(entries, catalogEntry)
	}

	return &Catalog{
		Streams: entries,
	}
}
