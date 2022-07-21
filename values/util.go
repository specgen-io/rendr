package values

func StrPtr(value string) *string {
	return &value
}

func BoolPtr(value bool) *bool {
	return &value
}
