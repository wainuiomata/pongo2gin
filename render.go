// Package pongo2gin is a template renderer that can be used with the Gin
// web framework https://github.com/gin-gonic/gin it uses the Pongo2 template
// library https://github.com/flosch/pongo2
package pongo2gin

import (
	"net/http"

	"github.com/flosch/pongo2/v5"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

// RenderOptions is used to configure the renderer.
type RenderOptions struct {
	TemplateDir string
	TemplateSet *pongo2.TemplateSet
	ContentType string
}

// Pongo2Render is a custom Gin template renderer using Pongo2.
type Pongo2Render struct {
	Options  *RenderOptions
	Template *pongo2.Template
	Context  pongo2.Context
}

// New creates a new Pongo2Render instance with custom Options.
func New(options RenderOptions) *Pongo2Render {
	// If TemplateSet is nil, rather than using pongo2.DefaultSet,
	// construct a new TemplateSet with the correct base directory.
	// This avoids the need to call pongo2.DefaultLoader.SetBaseDir,
	// and is necessary to support multiple Pongo2Render instances.
	if options.TemplateSet == nil {
		loader := pongo2.MustNewLocalFileSystemLoader(options.TemplateDir)
		options.TemplateSet = pongo2.NewSet(options.TemplateDir, loader)
	}

	return &Pongo2Render{
		Options: &options,
	}
}

// Default creates a Pongo2Render instance with default options.
func Default() *Pongo2Render {
	return New(RenderOptions{
		TemplateDir: "templates",
		TemplateSet: nil,
		ContentType: "text/html; charset=utf-8",
	})
}

// Instance should return a new Pongo2Render struct per request and prepare
// the template by either loading it from disk or using pongo2's cache.
func (p Pongo2Render) Instance(name string, data interface{}) render.Render {
	var template *pongo2.Template

	// Use template cache in Production mode.
	// In Debug mode load the file from disk each time.
	if gin.Mode() == "debug" {
		template = pongo2.Must(p.Options.TemplateSet.FromFile(name))
	} else {
		template = pongo2.Must(p.Options.TemplateSet.FromCache(name))
	}

	// This could be modifying the original data, need to check this.
	data.(pongo2.Context).Update(p.Options.TemplateSet.Globals)

	return Pongo2Render{
		Template: template,
		Context:  data.(pongo2.Context),
		Options:  p.Options,
	}
}

// Render should render the template to the response.
func (p Pongo2Render) Render(w http.ResponseWriter) error {
	p.WriteContentType(w)
	err := p.Template.ExecuteWriter(p.Context, w)
	return err
}

// WriteContentType should add the Content-Type header to the response
// when not set yet.
func (p Pongo2Render) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{p.Options.ContentType}
	}
}
