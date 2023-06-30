package datatypes

type AnyT interface{}

type StringT struct {
	Data       string
	Expiration int64
}

type ListT struct {
	Data       []string
	Expiration int64
}

type HashMapT struct {
	Data       map[string]string
	Expiration int64
}
