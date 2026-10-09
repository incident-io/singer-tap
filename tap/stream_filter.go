package tap

import (
	"context"
	"time"

	kitlog "github.com/go-kit/log"
	incident "github.com/incident-io/sdk-go"
	"github.com/incident-io/singer-tap/model"
)

type Filter struct {
	Stream       Stream
	CatalogEntry CatalogEntry
	// Since is the bookmark to fetch from, used only if Stream is an IncrementalStream.
	Since *time.Time
}

func (s *Filter) Output() *Output {
	output := s.Stream.Output()

	// // We need to filter the schema based on the catalog entry we have
	output.Schema.Properties = s.filterProperties(output.Schema.Properties, s.CatalogEntry)

	return output
}

func (s *Filter) GetRecords(ctx context.Context, logger kitlog.Logger, cl *incident.ClientWithResponses) ([]map[string]any, error) {
	var (
		records []map[string]any
		err     error
	)
	if incremental, ok := s.Stream.(IncrementalStream); ok {
		records, err = incremental.GetRecordsSince(ctx, logger, cl, s.Since)
	} else {
		records, err = s.Stream.GetRecords(ctx, logger, cl)
	}
	if err != nil {
		return nil, err
	}

	disabledFields := s.CatalogEntry.GetDisabledFields()

	// Filter out the disabled fields from each record (ew)
	for _, record := range records {
		for fieldName := range disabledFields {
			delete(record, fieldName)
		}
	}

	return records, nil
}

func (s *Filter) filterProperties(properties map[string]model.Property, catalogEntry CatalogEntry) map[string]model.Property {
	filteredProperties := map[string]model.Property{}
	disabledFields := catalogEntry.GetDisabledFields()

	for propertyName, property := range properties {
		// If the field is disabled, skip it
		if _, ok := disabledFields[propertyName]; ok {
			continue
		}

		filteredProperties[propertyName] = property
	}

	return filteredProperties
}
