package tap

import "github.com/incident-io/singer-tap/model"

type Metadata struct {
	// Pointer to where in the schmea this metadata applies
	Breadcrumb []string `json:"breadcrumb"`

	// Fields set for this metadata object
	Metadata MetadataFields `json:"metadata"`
}

type MetadataFields struct {
	/****
	* NON DISCOVERABLE FIELDS
	* We don't control these - pull them in and use them
	****/

	// Selected: if this node is selected by the user to be emitted
	// Can be field level or whole stream
	Selected *bool `json:"selected,omitempty"`

	// ReplicationMethod: the replication method to use
	// we ignored for our tap
	ReplicationMethod *string `json:"replicate-method,omitempty"`

	// ReplicationKey: the replicate key for this node
	// Used as a bookmark - ignore for our tap
	ReplicationKey *string `json:"replication-key,omitempty"`

	// ViewKeyProperties: not sure how this is used
	// ignored for our tap
	ViewKeyProperties *[]string `json:"view-key-properties,omitempty"`

	/****
	* DISCOVERABLE FIELDS
	* We can read and write these fields
	****/

	// Inclusion: whether we emit this field automatically
	// can be available (you choose), automatic (we choose), or unsupported (we don't emit)
	Inclusion string `json:"inclusion"`

	// SelectedByDefault: If the user doesn't specify should we
	// emit this field by default
	// This really only applies to available inclusion setting
	SelectedByDefault bool `json:"selected-by-default"`

	// ForcedReplicateMethod: we will set to FULL_TABLE for our tap
	ForcedReplicationMethod string `json:"forced-replication-method"`
}

func (m Metadata) DefaultMetadata(schema model.Schema) []Metadata {
	// By default we always include a top level metadata with the same
	// settings
	var metadata = []Metadata{
		{
			Breadcrumb: []string{},
			Metadata: MetadataFields{
				Inclusion:               "available",  // always set to available at stream level
				SelectedByDefault:       true,         // lets assume people always want our data
				ForcedReplicationMethod: "FULL_TABLE", // HIGHWAY TO THE DATA ZONE
			},
		},
	}

	return metadata
}
