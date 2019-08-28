package middleman

// middlemanError is returned if middleman listener falis to start
type middlemanError struct {
	errorString string
	innerError  error
}

type middlemanErrorInterface interface {
	Error() string

	// InnerError returns the inner error that occured
	InnerError() error
}

func (me middlemanError) InnerError() error {
	return me.innerError
}

func (me middlemanError) Error() string {
	return me.errorString
}

// ListenerError is returned when the listener failed to start
type ListenerError struct {
	middlemanError
}

// NewListenerError creates a new ListenerError
func NewListenerError(err error, errString string) ListenerError {
	return ListenerError{
		middlemanError: middlemanError{
			innerError:  err,
			errorString: errString,
		},
	}
}
