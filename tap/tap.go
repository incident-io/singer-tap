package tap

import (
	"context"
	"time"

	kitlog "github.com/go-kit/log"
	"github.com/incident-io/singer-tap/client"
)

func Sync(ctx context.Context, logger kitlog.Logger, ol *OutputLogger, cl *client.ClientWithResponses) error {
	for name, stream := range streams {
		logger := kitlog.With(logger, "stream", name)

		logger.Log("msg", "outputting schema")
		if err := ol.Log(stream.Output()); err != nil {
			return err
		}

		timeExtracted := time.Now().Format(time.RFC3339)
		logger.Log("msg", "loading records", "time_extracted", timeExtracted)
		records, err := stream.GetRecords(ctx, logger, cl)
		if err != nil {
			return err
		}

		logger.Log("msg", "outputting records", "count", len(records))
		for _, record := range records {
			op := &Output{
				Type:          OutputTypeRecord,
				Stream:        name,
				Record:        record,
				TimeExtracted: timeExtracted,
			}
			if err := ol.Log(op); err != nil {
				return err
			}
		}
	}

	return nil
}

func Discover(ctx context.Context, logger kitlog.Logger, ol *OutputLogger) error {
	catalog := NewCatalog(streams)

	if err := ol.CataLog(catalog); err != nil {
		return err
	}

	return nil
}
