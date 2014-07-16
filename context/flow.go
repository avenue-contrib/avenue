package context

import (
	"errors"
)

const (
	ErrUnhandled string = "request went unhandled"
)

func (c *Context) Next() {
	c.handlerIndex += 1

	// TODO: more rigid bounds checking
	if c.handlerIndex < len(c.Handlers) {
		c.Handlers[c.handlerIndex](c)
	} else {
		c.Error(500, errors.New(ErrUnhandled))
	}
}

func (c *Context) Abort(code int) {
	c.Response.WriteHeader(code)
}

func (c *Context) Fail(code int, err error) {
	c.Error(code, err)
	c.Abort(code)
}

func (c *Context) Error(code int, err error) {
	// c.errors = append(c.errors, err)
	c.Response.WriteHeader(code)
	c.Response.Write([]byte(err.Error()))
}

func (c *Context) LastError() error {
	s := len(c.errors)
	if s > 0 {
		return c.errors[s-1]
	} else {
		return nil
	}
}
