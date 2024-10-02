package errors

type FieldError struct {
	Value   interface{}
	Field   string
	Message string
}

func NewFieldError(v interface{}, f, m string) *FieldError {
	return &FieldError{Value: v, Field: m, Message: m}
}

func (v *FieldError) Error() string {
	return v.Message
}
