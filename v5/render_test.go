package pongo2gin

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flosch/pongo2/v5"
	"github.com/gin-gonic/gin"
)

// TestNew ensures that New creates a TemplateSet for options.TemplateDir
func TestNew(t *testing.T) {
	render := New(RenderOptions{
		TemplateDir: "templates/admin",
		TemplateSet: nil,
		ContentType: "text/html; charset=utf-8",
	})

	// Ensure a new TemplateSet was constructed.
	if render.Options.TemplateSet == nil {
		t.Fatal("Expect TemplateSet to not be nil")
	}

	// Check if the correct template is loaded.
	template, err := render.Options.TemplateSet.FromFile("admin.html")

	// Check that template loaded successfully.
	if err != nil {
		t.Fatal(err)
	}

	if template == nil {
		t.Fatal("Expect template to not be nil")
	}

	buffer := new(bytes.Buffer)
	err = template.ExecuteWriter(pongo2.Context{"country": "New Zealand"}, buffer)

	// Check that template rendered successfully.
	if err != nil {
		t.Fatal(err)
	}

	// Some quick checks to see if this is the template we are expecting.
	content := buffer.String()
	strings.Contains(content, "Admin Template")
	strings.Contains(content, "New Zealand")
}

// TestNew_TemplateSet tests the New function with a custom TemplateSet.
func TestNew_TemplateSet(t *testing.T) {
	name := "Email Template"
	expected := "pongo2gin@example.com"
	loader := pongo2.MustNewLocalFileSystemLoader("templates/email")

	render := New(RenderOptions{
		TemplateDir: "",
		TemplateSet: pongo2.NewSet("email", loader),
		ContentType: "text/html; charset=utf-8",
	})

	// Check that template loaded successfully.
	template, err := render.Options.TemplateSet.FromFile("email.html")
	if err != nil {
		t.Fatal(err)
	}

	if template == nil {
		t.Fatal("Expect template to not be nil")
	}

	// Check that template rendered successfully.
	buffer := new(bytes.Buffer)
	err = template.ExecuteWriter(pongo2.Context{"email": expected}, buffer)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the correct template is loaded.
	if !strings.Contains(buffer.String(), name) {
		t.Fatalf("Expected template %s", name)
	}

	// Check the content is what we expect.
	if !strings.Contains(buffer.String(), expected) {
		t.Fatalf("Expected text %s not found in body", expected)
	}
}

// TestNew_Debug ensures that the GIN_MODE=debug sets TemplateSet.Debug to true.
// This will mean the template cache is not used in Pongo2.
func TestNew_Debug(t *testing.T) {
	// Backup mode and restore for next test
	mode := gin.Mode()
	defer gin.SetMode(mode)

	// Set gin mode debug
	gin.SetMode("debug")

	render := New(RenderOptions{
		TemplateDir: "templates",
		TemplateSet: nil,
		ContentType: "text/html; charset=utf-8",
	})

	if render.Options.TemplateSet.Debug == false {
		t.Fatal("Expect TemplateSet.Debug to be true")
	}
}

// TestNew_Release ensures that the GIN_MODE=release sets TemplateSet.Debug to true.
func TestNew_Release(t *testing.T) {
	// Backup mode and restore for next test
	mode := gin.Mode()
	defer gin.SetMode(mode)

	// Set gin mode release
	gin.SetMode("release")

	render := New(RenderOptions{
		TemplateDir: "templates",
		TemplateSet: nil,
		ContentType: "text/html; charset=utf-8",
	})

	if render.Options.TemplateSet.Debug == true {
		t.Fatal("Expect TemplateSet.Debug to be false")
	}
}

// TestDefault tests the Default function which calls New with default options.
func TestDefault(t *testing.T) {
	name := "Index Template"
	expected := "Frodo"
	render := Default()
	template, err := render.Options.TemplateSet.FromFile("index.html")

	// Check that template loaded successfully.
	if err != nil {
		t.Fatal(err)
	}

	if template == nil {
		t.Fatal("Expect template to not be nil")
	}

	// Check that template rendered successfully.
	buffer := new(bytes.Buffer)
	err = template.ExecuteWriter(pongo2.Context{"name": expected}, buffer)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the correct template is loaded.
	if !strings.Contains(buffer.String(), name) {
		t.Fatalf("Expected template %s", name)
	}

	// Check the content is what we expect.
	if !strings.Contains(buffer.String(), expected) {
		t.Fatalf("Expected text %s not found in body", expected)
	}
}

// TestPongo2Render_Instance checks that the template is loaded properly.
func TestPongo2Render_Instance(t *testing.T) {
	render := Default()
	instance := render.Instance("index.html", pongo2.Context{"name": "Gandalf"})

	// Check that instance was created successfully.
	if instance == nil {
		t.Fatal("Instance returned should not be nil")
	}

	// Check that our Context vars carried over.
	if instance.(Pongo2Render).Context["name"] != "Gandalf" {
		t.Fatal("Expected context variable missing")
	}

	// Just checking Template is not nil should be enough.
	if instance.(Pongo2Render).Template == nil {
		t.Fatal("Expected template should not be nil")
	}
}

// TestPongo2Render_Render checks that the template renders to the response.
func TestPongo2Render_Render(t *testing.T) {
	expected := "Legolas"
	render := Default()
	instance := render.Instance("index.html", pongo2.Context{"name": expected})

	// Run the render method with request recorder.
	rr := httptest.NewRecorder()
	err := instance.Render(rr)
	if err != nil {
		t.Fatal(err)
	}

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the content is what we expect.
	if !strings.Contains(rr.Body.String(), expected) {
		t.Fatalf("Expected text %s not found in response body", expected)
	}
}
