package context

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

var (
	TemplatePath       string
	TemplateDelims     string
	TemplateExtensions []string
)

type Context struct {
	// handler and route variables
	StartTime  time.Time
	Parameters httprouter.Params

	// request and response wrappers
	Request  http.Request
	Response *ResponseWriter

	// context holdings
	Data     map[string]interface{}
	Handlers []Handler

	// bookkeeping
	handlerIndex int
	templatePath *string
}

// now that we have a context, define the handler type
type Handler func(ctx *Context)

func New(req http.Request, res http.ResponseWriter) *Context {
	return &Context{
		StartTime:    time.Now(),
		Request:      req,
		Response:     Wrap(res),
		Data:         make(map[string]interface{}),
		Handlers:     nil,
		handlerIndex: -1,
	}
}

func (c *Context) Param(key string) interface{} {
	return c.Parameters.ByName(key)
}

func (c *Context) Set(key string, item interface{}) {
	if c.Data == nil {
		c.Data = make(map[string]interface{})
	}
	c.Data[key] = item
}

func (c *Context) Get(key string) interface{} {
	if item, exists := c.Data[key]; exists {
		return item
	}

	return nil
}
