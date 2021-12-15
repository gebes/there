package there

import (
	"context"
	"testing"
	"time"
)

func TestServer_Start(t *testing.T) {
	router := NewRouter()
	go func() {
		err := router.Listen(8080)
		if err != nil {
			t.Error("unexpected error:", err)
		}
	}()
	router.shutdown(context.Background())
	time.Sleep(time.Millisecond * 100)
}

func TestServerTsl_Start(t *testing.T) {
	router := NewRouter()
	go func() {
		err := router.ListenToTLS(8081, "./test/server.crt", "./test/server.key")
		if err != nil {
			t.Error("unexpected error:", err)
		}
	}()
	router.shutdown(context.Background())
	time.Sleep(time.Millisecond * 100)
}
