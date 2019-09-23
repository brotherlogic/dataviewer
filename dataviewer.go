package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/brotherlogic/goserver"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbdc "github.com/brotherlogic/datacollector/proto"
	pbg "github.com/brotherlogic/goserver/proto"
)

type collector interface {
	GetDataSets(ctx context.Context, req *pbdc.GetDataSetsRequest) (*pbdc.GetDataSetsResponse, error)
}

type prodCollector struct {
	dial func(server string) (*grpc.ClientConn, error)
}

func (p *prodCollector) GetDataSets(ctx context.Context, req *pbdc.GetDataSetsRequest) (*pbdc.GetDataSetsResponse, error) {
	conn, err := p.dial("datacollector")
	if err != nil {
		return &pbdc.GetDataSetsResponse{}, err
	}
	defer conn.Close()

	client := pbdc.NewDataCollectorServiceClient(conn)
	return client.GetDataSets(ctx, req)
}

//Server main server type
type Server struct {
	*goserver.GoServer
	collector collector
}

// Init builds the server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		&prodCollector{},
	}
	s.collector = &prodCollector{dial: s.DialMaster}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	// Do nothing
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

type properties struct {
	Names []string
}

func (s *Server) deliver(w http.ResponseWriter, r *http.Request) {
	data, err := Asset("templates/main.html")
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error: %v", err))
		return
	}
	sets, err := s.getDataSets(context.Background())
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error: %v", err))
		return
	}
	props := properties{Names: sets}
	err = s.render(string(data), props, w)
	if err != nil {
		fmt.Fprintf(w, fmt.Sprintf("Error: %v", err))
	}
}

func (s *Server) serveUp() {
	http.HandleFunc("/", s.deliver)
	err := http.ListenAndServe(":8086", nil)
	if err != nil {
		panic(err)
	}
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server
	err := server.RegisterServer("dataviewer", false)
	if err != nil {
		log.Fatalf("Unable to register: %v", err)
	}
	go server.serveUp()
	fmt.Printf("%v", server.Serve())
}
