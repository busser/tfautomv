package rules

type baseRule struct {
	resourceType string
	attribute    string
}

func (r baseRule) AppliesTo(resourceType, attribute string) bool {
	return resourceType == r.resourceType && attribute == r.attribute
}
