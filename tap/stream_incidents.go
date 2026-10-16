package tap

import (
	"context"
	"time"

	kitlog "github.com/go-kit/log"
	incident "github.com/incident-io/sdk-go"
	"github.com/incident-io/singer-tap/model"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func init() {
	register(&StreamIncidents{})
}

type StreamIncidents struct {
}

func (s *StreamIncidents) Output() *Output {
	return &Output{
		Type:   OutputTypeSchema,
		Stream: "incidents",
		Schema: &model.Schema{
			HasAdditionalProperties: false,
			Type:                    []string{"object"},
			Properties:              model.IncidentV2.Schema().Properties,
		},
		KeyProperties:      []string{"id"},
		BookmarkProperties: []string{BookmarkKey},
	}
}

func (s *StreamIncidents) GetRecords(ctx context.Context, logger kitlog.Logger, cl *incident.ClientWithResponses) ([]map[string]any, error) {
	return s.GetRecordsSince(ctx, logger, cl, nil)
}

func (s *StreamIncidents) GetRecordsSince(ctx context.Context, logger kitlog.Logger, cl *incident.ClientWithResponses, since *time.Time) ([]map[string]any, error) {
	var (
		after    *string
		pageSize = int64(250)
		results  = []map[string]any{}
		seen     = map[string]bool{}
	)

	for {
		logger.Log("msg", "loading page", "page_size", pageSize, "after", after, "since", since)
		params := &incident.IncidentsV2ListParams{
			PageSize: &pageSize,
			After:    after,
		}
		if since != nil {
			params.UpdatedAt = UpdatedSinceFilter(*since)
		}

		page, err := cl.IncidentsV2ListWithResponse(ctx, params)
		if err != nil {
			return nil, errors.Wrap(err, "listing incidents")
		}

		for _, element := range page.JSON200.Incidents {
			// If the incident we passed as after is deleted or changes mode mid-sync, the
			// endpoint can't find it and starts again from the first page. Skip what we've
			// already loaded and keep paging, which still ends on an empty page.
			if seen[element.Id] {
				continue
			}
			seen[element.Id] = true

			attachments, err := s.GetAttachments(ctx, logger, cl, element.Id)
			if err != nil {
				return nil, errors.Wrap(err, "listing incident attachments")
			}

			updates, err := s.GetIncidentUpdates(ctx, logger, cl, element.Id)
			if err != nil {
				return nil, errors.Wrap(err, "listing incident updates")
			}

			results = append(results, model.IncidentV2.Serialize(element, attachments, updates))
		}
		if count := len(page.JSON200.Incidents); count == 0 {
			return results, nil // end pagination
		} else {
			after = lo.ToPtr(page.JSON200.Incidents[count-1].Id)
		}
	}
}

func (s *StreamIncidents) GetAttachments(ctx context.Context, logger kitlog.Logger, cl *incident.ClientWithResponses, incidentId string) ([]incident.IncidentAttachmentV1, error) {
	var (
		results = []incident.IncidentAttachmentV1{}
	)

	response, err := cl.IncidentAttachmentsV1ListWithResponse(ctx, &incident.IncidentAttachmentsV1ListParams{
		IncidentId: &incidentId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "listing attachments for incidents stream")
	}

	results = append(results, response.JSON200.IncidentAttachments...)
	return results, nil
}

func (s *StreamIncidents) GetIncidentUpdates(ctx context.Context, logger kitlog.Logger, cl *incident.ClientWithResponses, incidentId string) ([]incident.IncidentUpdateV2, error) {
	var (
		after    *string
		pageSize int64 = 250
		results        = []incident.IncidentUpdateV2{}
	)

	for {
		logger.Log("msg", "loading incident updates page", "page_size", pageSize, "after", after)
		page, err := cl.IncidentUpdatesV2ListWithResponse(ctx, &incident.IncidentUpdatesV2ListParams{
			IncidentId: &incidentId,
			PageSize:   &pageSize,
			After:      after,
		})
		if err != nil {
			return nil, errors.Wrap(err, "listing incident updates")
		}

		results = append(results, page.JSON200.IncidentUpdates...)

		if count := len(page.JSON200.IncidentUpdates); count == 0 {
			return results, nil // end pagination
		} else {
			after = lo.ToPtr(page.JSON200.IncidentUpdates[count-1].Id)
		}
	}
}
