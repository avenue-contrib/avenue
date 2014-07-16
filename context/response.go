package context

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"errors"
	"html/template"
	"net"
	"net/http"
	"path/filepath"
)

const (
	StatusUnset int = -1
)

//-----------------------------------------------
// ResponseWriter Wrapper
//-----------------------------------------------

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func Wrap(res http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{res, StatusUnset}
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ResponseWriter) Write(data []byte) (int, error) {
	// net/http.Response.Write only has two options: 200 or 500
	// we will follow that lead and defer to their logic

	// check if the write gave an error and set status accordingly
	size, err := w.ResponseWriter.Write(data)
	if err != nil {
		// error on write, we give a 500
		w.status = http.StatusInternalServerError
	} else if w.WasWritten() == false {
		// everything went okay and we never set a custom
		// status so 200 it is
		w.status = http.StatusOK
	}

	// can easily tap into Content-Length here with 'size'
	return size, err
}

// returns the status of the given response
func (w *ResponseWriter) Status() int {
	return w.status
}

// return a boolean acknowledging if a status code has all ready been set
func (w *ResponseWriter) WasWritten() bool {
	return w.status == StatusUnset
}

// allow connection hijacking
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}
	return hijacker.Hijack()
}

//-----------------------------------------------
// Context Controls
//-----------------------------------------------

func (c *Context) JSON(code int, obj interface{}) {
	c.Response.Header().Set("Content-Type", "application/json")
	if code >= 0 {
		c.Response.WriteHeader(code)
	}
	encoder := json.NewEncoder(c.Response)
	if err := encoder.Encode(obj); err != nil {
		c.Error(500, err)
		http.Error(c.Response, err.Error(), 500)
	}
}

func (c *Context) XML(code int, obj interface{}) {
	c.Response.Header().Set("Content-Type", "application/xml")
	if code >= 0 {
		c.Response.WriteHeader(code)
	}
	encoder := xml.NewEncoder(c.Response)
	if err := encoder.Encode(obj); err != nil {
		c.Error(500, err)
		http.Error(c.Response, err.Error(), 500)
	}
}

func (c *Context) Render(code int, name string, data interface{}) {
	templ, err := template.ParseFiles(filepath.Join(TemplatePath, name))
	if err != nil {
		// TODO: check for prod environment
		c.Error(500, err)
		return
	}

	err = templ.Execute(c.Response, data)
	if err != nil {
		// TODO: check for prod environment
		c.Error(500, err)
		return
	} else {
		c.Response.Header().Set("Content-Type", "text/html")
	}
}

func (c *Context) String(code int, msg string) {
	if code >= 0 {
		c.Response.WriteHeader(code)
	}
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.Write([]byte(msg))
}

func (c *Context) Bytes(code int, data []byte) {
	c.Response.WriteHeader(code)
	c.Response.Write(data)
}
