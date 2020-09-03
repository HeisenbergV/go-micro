package server

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/util/http"
)

func TestHTTPServer(t *testing.T) {
	srv := NewServer(Address("localhost:8000"))
	mux := http.NewHttpServ()
	mux.GET("/aa", func(content *gin.Context) {
		content.JSON(200, "world")
	})

	hd := srv.NewHandler(mux)
	srv.Handle(hd)
	// start server
	if err := srv.Start(); err != nil {
		t.Fatal(err)
	}

	for {

	}
}
