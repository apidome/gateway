package middleman

import "net/http"

// Get Adds a GET middleware to a route
// the 'path' argument will be prefixed with a '^' and
// suffixed with a '$' for regex matching
func (mm *Middleman) Get(path string, middleware Middleware) error {
	return mm.addMiddleware(path, http.MethodGet, middleware)
}

// Head Adds a HEAD middleware to a route
// the 'path' argument will be prefixed with a '^' and
// suffixed with a '$' for regex matching
func (mm *Middleman) Head(path string, middleware Middleware) error {
	return mm.addMiddleware(path, http.MethodHead, middleware)
}

// Post Adds a POST middleware to a route
// the 'path' argument will be prefixed with a '^' and
// suffixed with a '$' for regex matching
func (mm *Middleman) Post(path string, middleware Middleware) error {
	return mm.addMiddleware(path, http.MethodPost, middleware)
}

// Put Adds a PUT middleware to a route
// the 'path' argument will be prefixed with a '^' and
// suffixed with a '$' for regex matching
func (mm *Middleman) Put(path string, middleware Middleware) error {
	return mm.addMiddleware(path, http.MethodPut, middleware)
}

// Delete Adds a DELETE middleware to a route
// the 'path' argument will be prefixed with a '^' and
// suffixed with a '$' for regex matching
func (mm *Middleman) Delete(path string, middleware Middleware) error {
	return mm.addMiddleware(path, http.MethodDelete, middleware)
}

// Connect Adds a CONNECT middleware to a route
// the 'path' argument will be prefixed with a '^' and
// suffixed with a '$' for regex matching
func (mm *Middleman) Connect(path string, middleware Middleware) error {
	return mm.addMiddleware(path, http.MethodConnect, middleware)
}

// Options Adds a OPTIONS middleware to a route
// the 'path' argument will be prefixed with a '^' and
// suffixed with a '$' for regex matching
func (mm *Middleman) Options(path string, middleware Middleware) error {
	return mm.addMiddleware(path, http.MethodOptions, middleware)
}

// Trace Adds a TRACE middleware to a route
// the 'path' argument will be prefixed with a '^' and
// suffixed with a '$' for regex matching
func (mm *Middleman) Trace(path string, middleware Middleware) error {
	return mm.addMiddleware(path, http.MethodTrace, middleware)
}

// Patch Adds a PATCH middleware to a route
// the 'path' argument will be prefixed with a '^' and
// suffixed with a '$' for regex matching
func (mm *Middleman) Patch(path string, middleware Middleware) error {
	return mm.addMiddleware(path, http.MethodPatch, middleware)
}

// All Adds a middleware to all methods of a route
// the 'path' argument will be prefixed with a '^' and
// suffixed with a '$' for regex matching
func (mm *Middleman) All(path string, middleware Middleware) error {
	var err error
	for _, method := range methods {
		err = mm.addMiddleware(path, method, middleware)

		if err != nil {
			return err
		}
	}

	return nil
}
