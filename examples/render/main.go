package main

import (
	"avenue"
	"avenue/context"
)

func main() {
	serv := avenue.New()

	serv.Templates("templates", "{{ }}", []string{"html", "tmpl"})

	serv.GET("/:name", func(ctx *context.Context) {
		// use .Param(key string) function to lookup the name
		// have to assert type
		path := ctx.Param("name").(string)

		ctx.Render(200, "index.html", map[string]interface{}{
			"Endpoint": path,
		})
	})

	serv.GET("/:name/tmpl", func(ctx *context.Context) {
		// use .Param(key string) function to lookup the name
		// have to assert type
		path := ctx.Param("name").(string)

		ctx.Render(200, "index.tmpl", map[string]interface{}{
			"Endpoint": path,
		})
	})

	serv.Run(":8080")
}
