package main

import (
	"fmt"

	"github.com/avenue-contrib/avenue"
	"github.com/avenue-contrib/avenue/context"
)

func main() {
	serv := avenue.New()
	serv.GET("/:name", func(ctx *context.Context) {
		// use .Param(key string) function to lookup the name
		// have to assert type
		path := ctx.Param("name").(string)

		// send back a string with status 200
		ctx.String(200, fmt.Sprintf("Hello, %s!", path[1:]))
	})
	serv.Run(":8080")
}
