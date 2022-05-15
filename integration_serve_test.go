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
		if err.Error() != "http: Server closed" {
			t.Error("unexpected error:", err)
		}
	}()
	time.Sleep(time.Millisecond * 10)
	err := router.Server.Shutdown(context.Background())
	if err != nil {
		t.Error("unexpected error:", err)
	}
	time.Sleep(time.Millisecond * 50)
}

func TestServerTsl_Start(t *testing.T) {
	router := NewRouter()
	go func() {
		err := router.ListenToTLS(8081, "./test/server.crt", "./test/server.key")
		if err.Error() != "http: Server closed" {
			t.Error("unexpected error:", err)
		}
	}()
	time.Sleep(time.Millisecond * 10)
	err := router.Server.Shutdown(context.Background())
	if err != nil {
		t.Error("unexpected error:", err)
	}
	time.Sleep(time.Millisecond * 50)
}
