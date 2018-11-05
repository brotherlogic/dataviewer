package main

import (
	"testing"
)

func InitTestServer() *Server {
	s := Init()
	return s
}

func TestDoNothing(t *testing.T) {
	doNothing()
}
