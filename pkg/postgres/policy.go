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

// CreatePolicyRequest is a struct used to represent info about the whole policy. It
// contains a public key, a label, and a list of operations to be applied to the
// data if the user applies this policy to their event stream.
type CreatePolicyRequest struct {
	PublicKey  string
	Label      string
	Operations []Operation
}

// CreatePolicyResponse is returned from the Postgres module when a policy
// record is successfully inserted. THe RPC layer will convert this type into
// the external type we serialize over the wire.
type CreatePolicyResponse struct {
	ID    string
	Token string
	CreatePolicyRequest
}

// DeletePolicyRequest is a type used to pass incoming delete requests from the
// RPC handler function to the Postgres DB layer.
type DeletePolicyRequest struct {
	ID    string
	Token string
}

// PolicyResponse is a type used to return a list of policies from postgoreos
// that should be returned to the client. Each record contains the ID of the
// policy, its label, list of operations and the associated public key.
type PolicyResponse struct {
	ID         string
	Label      string
	Operations []Operation
	PublicKey  string
}
