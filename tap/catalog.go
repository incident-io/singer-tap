package tap

import (
	"github.com/incident-io/singer-tap/model"
)

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
	Metadata *[]Metadata `json:"metadata,omitempty"`
}

// Actual catalog that we export
// contains an array of all our streams
type Catalog struct {
	Streams []CatalogEntry `json:"streams"`
}

func (c *Catalog) GetEnabledStreams() []CatalogEntry {
	var enabledStreams []CatalogEntry

	// Go through all streams registered in the catalog
	for _, entry := range c.Streams {
		// if there is no metadata then just include the stream
		if entry.Metadata == nil {
			enabledStreams = append(enabledStreams, entry)
		} else {
			for _, metadata := range *entry.Metadata {
				// Only check the top level metadata
				if len(metadata.Breadcrumb) > 0 {
					continue
				}

				// Check if the metadata has the user input "selected" bool
				if metadata.Metadata.Selected != nil {
					// If so, check its set to true!
					if *metadata.Metadata.Selected {
						enabledStreams = append(enabledStreams, entry)
					}
					// otherwise check if WE have set to select this by default
				} else if metadata.Metadata.SelectedByDefault {
					enabledStreams = append(enabledStreams, entry)
				}
			}
		}
	}

	return enabledStreams
}

func (c *CatalogEntry) GetDisabledFields() map[string]bool {
	// Just something to enable quick lookups of fields by name
	var disabledFields = map[string]bool{}

	// For the given stream, get the enabled fields
	// For this catalog entry, get the metadata, and build a list of all the enabled fields
	for _, metadata := range *c.Metadata {
		// Ignore the top level metadata
		if len(metadata.Breadcrumb) == 0 {
			continue
		}

		// Check if the metadata has the user input "selected" bool
		if metadata.Metadata.Selected != nil {
			// If so, check its set to false!
			if !*metadata.Metadata.Selected {
				disabledFields[metadata.Breadcrumb[len(metadata.Breadcrumb)-1]] = true
			}
		} else {
			// There's no selected key, so check if WE have set the selected by default
			if !metadata.Metadata.SelectedByDefault {
				disabledFields[metadata.Breadcrumb[len(metadata.Breadcrumb)-1]] = true
			}
		}
	}

	return disabledFields
}

func NewDefaultCatalog(streams map[string]Stream) *Catalog {
	entries := []CatalogEntry{}

	for name, stream := range streams {
		streamSchema := *stream.Output().Schema
		metadata := Metadata{}.DefaultMetadata(streamSchema)

		catalogEntry := CatalogEntry{
			Stream:      name,
			TapStreamID: name,
			Schema:      streamSchema,
			Metadata:    &metadata,
		}

		entries = append(entries, catalogEntry)
	}

	return &Catalog{
		Streams: entries,
	}
}
