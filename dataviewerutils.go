package main

import (
	"html/template"
	"io"
)

func (s *Server) render(f string, w io.Writer) error {
	templ := template.New("main")
	templ, err := templ.Parse(f)
	if err != nil {
		return err
	}

	templ.Execute(w, "")
	return nil
}
