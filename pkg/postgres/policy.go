package postgres

// Action is a type alias for int - used when defining an iota based enum const
type Action int

const (
	// Share defines an action for an operation where the data is shared
	Share Action = iota + 1
	// Bin defines an action for an operation where the data is binned
	Bin
	// Average defines an action for an operation where a moving average is applied
	// to the data.
	Average
)

// String is an implementation of stringer for our custom type used in the
// iota/enum
func (a Action) String() string {
	names := [...]string{
		"Unknown",
		"Share",
		"Bin",
		"Moving Average",
	}

	if a < Share || a > Average {
		return names[0]
	}

	return names[a]
}

// Operation is a single transformative operation that should be applied to the
// incoming event stream.
type Operation struct {
	SensorID int       `json:"sensorID"`
	Action   Action    `json:"action"`
	Bins     []float64 `json:"bins"`
	Interval int       `json:"interval"`
}

// PolicyRequest is a struct used to represent info about the whole policy. It
// contains a public key, a label, and a list of operations to be applied to the
// data if the user applies this policy to their event stream.
type PolicyRequest struct {
	PublicKey  string      `db:"public_key"`
	Label      string      `db:"label"`
	Operations []Operation `db:"operations"`
}

// PolicyResponse is used...
type PolicyResponse struct {
	ID    int64  `db:"id"`
	Token string `db:"id"`
	PolicyRequest
}
