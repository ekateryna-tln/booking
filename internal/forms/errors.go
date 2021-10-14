package forms

type errors map[string][]string

// Add adds an error message for given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns first error message
func (e errors) Get(field string) string {
	if len(e[field]) == 0 {
		return ""
	}
	return e[field][0]
}
