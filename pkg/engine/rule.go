package engine

// A Rule allows tfautomv to equate certains attribute values that would
// normally be considered different.
type Rule interface {
	// Each rule has a unique string representation. This representation is how
	// users provides rules to tfautomv.
	String() string

	// Whether the rule applies to the given resource type and attribute.
	AppliesTo(resourceType, attribute string) bool

	// Whether the rule equates the two values.
	Equates(a, b interface{}) bool
}
