package main

import (
	"context"
	"errors"
	"fmt"
	datamanager "main/internal/dataManager"
	"main/internal/mensa"
	"main/internal/rvv"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	mensa.Init()
	rvv.Init()
	err := datamanager.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = http.ListenAndServe(":8123", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	shutdownCtx, shutdown := context.WithCancel(context.Background())
	defer shutdown()
	err = waitForInterrupt(shutdownCtx)
	if err != nil {
		fmt.Println(err)
	}
}

func waitForInterrupt(ctx context.Context) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-c:
		return fmt.Errorf("received signal %s", sig)
	case <-ctx.Done():
		return errors.New("canceled")
	}
}
