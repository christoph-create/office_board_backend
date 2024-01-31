package main

import (
	"context"
	"errors"
	"fmt"
	"main/internal/mensa"
	"main/internal/rvv"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	mensa mensa.Mensa
	rvv   rvv.Rvv
}

func main() {
	_, err := mensa.New()
	if err != nil {
		fmt.Println(err)
	}

	_, err = rvv.New()
	if err != nil {
		fmt.Println(err)
	}

	shutdownCtx, shutdown := context.WithCancel(context.Background())
	defer shutdown()
	err = waitForInterrupt(shutdownCtx)
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
