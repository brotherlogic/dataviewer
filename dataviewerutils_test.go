package main

import (
	"bytes"
	"testing"
)

func InitTestServer() *Server {
	s := Init()
	return s
}

func TestTemplateFailure(t *testing.T) {
	s := InitTestServer()
	var buf bytes.Buffer
	err := s.render("{{.broken", &buf)

	if err == nil {
		t.Errorf("No error in processing")
	}
}

func TestEasyTemplate(t *testing.T) {
	s := InitTestServer()
	var buf bytes.Buffer
	err := s.render("templates/main.html", &buf)

	if err != nil {
		t.Errorf("Rendering error: %v", err)
	}

	if len(buf.String()) == 0 {
		t.Errorf("Error in building string: %v", buf.String())
	}
}
