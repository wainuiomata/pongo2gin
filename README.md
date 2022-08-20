Pongo2gin
=========

Package pongo2gin is a template renderer that can be used with the Gin web
framework https://github.com/gin-gonic/gin it uses the Pongo2 template library
https://github.com/flosch/pongo2

This simple binding library is based on a similar library built for using
Handlebars templates with Gin: https://gitlab.com/go-box/ginraymond.

Requirements
------------

Requires Gin 1.16 or higher and Pongo2 version 5.

## Here are versions compatible with other versions of Pongo2

 [pongo2 version 5](https://gitlab.com/go-box/pongo2gin/) -  Compatible with pongo2 version 5
 
 [pongo2 version 1](https://gitlab.com/go-box/pongo2gin/tree/main/v1) - Compatible with pongo2 version 1
 
 [pongo2 version 4](https://gitlab.com/go-box/pongo2gin/tree/main/v4) -  Compatible with pongo2 version 4

# Please don't forget to give stars :)

## Installation

`go get "gitlab.com/go-box/pongo2gin"`

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
import (
	"log"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"

	"gitlab.com/go-box/pongo2gin"
)

// GetAllData returns all posts
func GetAllData(c *gin.Context) {
	posts := []string{
		"Rob van der Linde",
		"John Curley",
		"Andrejs Cainikovs",
	}

	// Call the HTML method of the Context to render a template
	c.HTML(http.StatusOK, "index.html",
		pongo2.Context{
			"title": "hello pongo",
			"posts": posts,
		},
	)
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.Use(gin.Recovery())
	r.HTMLRender = pongo2gin.Default()
	r.GET("/", GetAllData)
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

If he Gin Mode is set to "debug" then templates will be loaded from disk on
each request.

Caching is implemented by the Pongo2 library itself.

GoDoc
-----

https://godoc.org/gitlab.com/go-box/pongo2gin
