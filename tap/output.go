package tap

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/incident-io/singer-tap/model"
)

// OutputType is the type of Singer tap output for each message.
type OutputType string

var (
	OutputTypeSchema OutputType = "SCHEMA"
	OutputTypeRecord OutputType = "RECORD"
)

// Output is what we log to STDOUT as a message provided to the downstream Singer target.
//
// This tap supports type types of output:
//
// - SCHEMA: Specifies the schema of this stream in JSON schema format.
// - RECORD: A record from the stream.
//
// We (currently) do not support the other types of output such as STATE.
type Output struct {
	// Type is the type of the stream, e.g. "SCHEMA" or "RECORD"
	Type OutputType `json:"type,omitempty"`
	// Stream is the name of the stream, e.g. "users"
	Stream string `json:"stream,omitempty"`
	// Schema is the schema of the stream, if Type == "SCHEMA", in JSON schema format.
	Schema *model.Schema `json:"schema,omitempty"`
	// Record is a record from the stream, if Type == "RECORD".
	Record map[string]any `json:"record,omitempty"`
	// TimeExtracted is the time that this record was extracted, if Type == "RECORD".
	TimeExtracted string `json:"time_extracted,omitempty"`
	// KeyProperties is a list of strings indicating which properties make up the primary
	// key for this stream.
	//
	// Each item in the list must be the name of a top-level property defined in the schema.
	// A value for key_properties must be provided, but it may be an empty list to indicate
	// that there is no primary key.
	KeyProperties []string `json:"key_properties,omitempty"`
	// BookmarkProperties is an optional list of strings indicating which properties
	// should be used to bookmark the stream, such as "last_updated_at".
	BookmarkProperties []string `json:"bookmark_properties,omitempty"`
}

// OutputLogger is a logger that logs to STDOUT in the format expected by the downstream
// Singer target.
type OutputLogger struct {
	w io.Writer
}

func NewOutputLogger(w io.Writer) *OutputLogger {
	return &OutputLogger{w: w}
}

func (o *OutputLogger) Log(op *Output) error {
	data, err := json.Marshal(op)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(o.w, string(data))
	if err != nil {
		return err
	}

	return nil
}

func (o *OutputLogger) CataLog(catalog *Catalog) error {
	data, err := json.Marshal(catalog)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(o.w, string(data))
	if err != nil {
		return err
	}

	return nil
}
