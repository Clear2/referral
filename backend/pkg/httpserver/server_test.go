package httpserver

import (
	"net"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestStartReturnsBindError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve port: %v", err)
	}
	defer listener.Close()
	server := New(gin.New(), func(server *serverImpl) { server.Address = listener.Addr().String() })
	if err := server.Start(); err == nil {
		t.Fatal("Start() error = nil, want bind error")
	}
}
