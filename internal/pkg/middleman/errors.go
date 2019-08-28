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
func NewListenerError(innerError error, errString string) error {
	return ListenerError{
		middlemanError: middlemanError{
			innerError:  innerError,
			errorString: errString,
		},
	}
}

// RegexCompilationError is returned when the path of a middleware
// failed regex compilation
type RegexCompilationError struct {
	middlemanError
}

// NewRegexCompilationError creates a new RegexCompilationError
func NewRegexCompilationError(innerError error,
	errString string) error {
	return RegexCompilationError{
		middlemanError: middlemanError{
			innerError:  innerError,
			errorString: errString,
		},
	}
}
