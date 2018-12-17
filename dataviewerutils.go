package main

import (
	"html/template"
	"io"

	"golang.org/x/net/context"

	pbdc "github.com/brotherlogic/datacollector/proto"
)

func (s *Server) getDataSets(ctx context.Context) ([]string, error) {
	resp, err := s.collector.GetDataSets(ctx, &pbdc.GetDataSetsRequest{})
	if err != nil {
		return []string{}, err
	}

	result := []string{}
	for _, res := range resp.DataSets {
		result = append(result, res.SpecName)
	}

	return result, nil
}

func (s *Server) render(f string, props properties, w io.Writer) error {
	templ := template.New("main")
	templ, err := templ.Parse(f)
	if err != nil {
		return err
	}

	templ.Execute(w, props)
	return nil
}
