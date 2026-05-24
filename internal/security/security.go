package security

// Decision expresses the result of a security evaluation.
type Decision string

const (
	DecisionAllow  Decision = "allow"
	DecisionDeny   Decision = "deny"
	DecisionReview Decision = "review"
)

// Subject identifies a caller or workload.
type Subject struct {
	Name   string
	Roles  []string
	Scopes []string
}

// Rule is a minimal policy scaffold used for future authorization checks.
type Rule struct {
	Name   string
	Action string
	Object string
	Effect Decision
}

// Evaluation is the result of applying a rule set to a subject.
type Evaluation struct {
	Subject Subject
	Rule    Rule
	Result  Decision
	Reason  string
}

// IsAllowed returns true when the decision is to allow.
func IsAllowed(decision Decision) bool {
	return decision == DecisionAllow
}
