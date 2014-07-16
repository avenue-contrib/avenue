package logging

import (
	"log"
	"time"

	"avenue/context"

	"github.com/mgutz/ansi"
)

var (
	red    = ansi.ColorCode("red+h:black")
	green  = ansi.ColorCode("green+h:black")
	yellow = ansi.ColorCode("yellow+h:black")
	reset  = ansi.ColorCode("reset")
)

func Plain() context.Handler {
	return func(ctx *context.Context) {
		// save the start time to do latency calculation
		start := time.Now()

		// save the IP of the requester
		requester := ctx.Request.Header.Get("X-Real-IP")

		// if the requester-header is empty, check the forwarded-header
		if requester == "" {
			requester = ctx.Request.Header.Get("X-Forwarded-For")
		}

		// if the requester is still empty, use the hard-coded address from the socket
		if requester == "" {
			requester = ctx.Request.RemoteAddr
		}

		// ... finally, log the fact we got a request
		log.Printf("<-- %16s | %6s | %s\n", requester, ctx.Request.Method, ctx.Request.URL.Path)

		// keep going
		ctx.Next()
		// and wait until we come back

		log.Printf("--> %16s | %6d | %s | %s\n",
			requester, ctx.Response.Status(), time.Since(start), ctx.Request.URL.Path,
		)
	}
}

func Color() context.Handler {
	return func(ctx *context.Context) {
		// save the start time to do latency calculation
		start := time.Now()

		// save the IP of the requester
		requester := ctx.Request.Header.Get("X-Real-IP")

		// if the requester-header is empty, check the forwarded-header
		if requester == "" {
			requester = ctx.Request.Header.Get("X-Forwarded-For")
		}

		// if the requester is still empty, use the hard-coded address from the socket
		if requester == "" {
			requester = ctx.Request.RemoteAddr
		}

		// ... finally, log the fact we got a request
		log.Printf("<-- %16s | %6s | %s\n", requester, ctx.Request.Method, ctx.Request.URL.Path)

		// keep going
		ctx.Next()
		// and wait until we come back

		var color string
		if code := ctx.Response.Status(); code == 200 {
			color = green
		} else if code >= 300 && code <= 399 {
			color = yellow
		} else {
			color = red
		}

		log.Printf("--> %16s | %s%6d%s | %s | %s\n",
			requester,
			color, ctx.Response.Status(), reset,
			time.Since(start), ctx.Request.URL.Path,
		)
	}
}
