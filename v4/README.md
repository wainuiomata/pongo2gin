Pongo2gin
=========

Package pongo2gin is a template renderer that can be used with the Gin web
framework https://github.com/gin-gonic/gin it uses the Pongo2 template library
https://github.com/flosch/pongo2

This simple binding library is based on a similar library built for using
Handlebars templates with Gin: https://gitlab.com/go-box/ginraymond.

Requirements
------------

Requires Go 1.18 or higher, Gin and Pongo2 version 4.

## Installation

`go get gitlab.com/go-box/pongo2gin/v4`

Usage
-----

To use pongo2gin you need to set your router.HTMLRenderer to a new renderer
instance, this is done after creating the Gin router when the Gin application
starts up. You can use pongo2gin.Default() to create a new renderer with
default options, this assumes templates will be located in the "templates"
directory, or you can use pongo2.New() to specify a custom location.

To render templates from a route, call c.HTML just as you would with
regular Gin templates, the only difference is that you pass template
data as a pongo2.Context instead of gin.H type.

Basic Example
-------------

```go
package main

import (
	"log"
	"net/http"

	"github.com/flosch/pongo2/v4"
	"github.com/gin-gonic/gin"

	"gitlab.com/go-box/pongo2gin/v4"
)

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.HTMLRender = pongo2gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", pongo2.Context{
			"name": "Pongo2",
		})
	})
	log.Fatal(r.Run(":8888"))
}

```

RenderOptions
-------------

When calling pongo2gin.New() instead of pongo2gin.Default() you can use these
custom RenderOptions:

```go
type RenderOptions struct {
    TemplateDir string              // location of the template directory
    TemplateSet *pongo2.TemplateSet // pongo2 template set with custom loader, or nil
    ContentType string              // Content-Type header used when calling c.HTML()
}
```

Template Caching
----------------

Templates will be cached if the current Gin Mode is set to anything but "debug",
this means the first time a template is used it will still load from disk, but
after that the cached template will be used from memory instead.

If the Gin Mode is set to "debug" then templates will be loaded from disk on
each request.

Caching is implemented by the Pongo2 library itself.

GoDoc
-----

https://godoc.org/gitlab.com/go-box/pongo2gin
