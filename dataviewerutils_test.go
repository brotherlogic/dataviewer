package main

import (
	"bytes"
	"fmt"
	"testing"

	pbdc "github.com/brotherlogic/datacollector/proto"
	"golang.org/x/net/context"
)

type testCollector struct {
	fail bool
}

func (t *testCollector) GetDataSets(ctx context.Context, req *pbdc.GetDataSetsRequest) (*pbdc.GetDataSetsResponse, error) {
	if t.fail {
		return &pbdc.GetDataSetsResponse{}, fmt.Errorf("Built to fail")
	}
	return &pbdc.GetDataSetsResponse{
		DataSets: []*pbdc.DataSet{
			&pbdc.DataSet{SpecName: "testing"},
		},
	}, nil
}
func InitTestServer() *Server {
	s := Init()
	s.collector = &testCollector{}
	return s
}

func TestGetDataSets(t *testing.T) {
	s := InitTestServer()
	sets, err := s.getDataSets(context.Background())
	if err != nil {
		t.Fatalf("Get Data sets has failed: %v", err)
	}

	if len(sets) != 1 {
		t.Errorf("Bad return on get data sets: %v", sets)
	}
}

func TestGetDataSetsOnFail(t *testing.T) {
	s := InitTestServer()
	s.collector = &testCollector{fail: true}
	sets, err := s.getDataSets(context.Background())
	if err == nil {
		t.Fatalf("Get Data sets has not failed: %v", sets)
	}

}

func TestTemplateFailure(t *testing.T) {
	s := InitTestServer()
	var buf bytes.Buffer
	err := s.render("{{.broken", properties{}, &buf)

	if err == nil {
		t.Errorf("No error in processing")
	}
}

func TestEasyTemplate(t *testing.T) {
	s := InitTestServer()
	var buf bytes.Buffer
	err := s.render("templates/main.html", properties{}, &buf)

	if err != nil {
		t.Errorf("Rendering error: %v", err)
	}

	if len(buf.String()) == 0 {
		t.Errorf("Error in building string: %v", buf.String())
	}
}
