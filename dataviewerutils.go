package main

import (
	"html/template"
	"io"
)

func (s *Server) render(f string, w io.Writer) error {
	templ, err := template.ParseFiles(f)
	if err != nil {
		return err
	}

	templ.Execute(w, "")
	return nil
}
