package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

// server： 提供http服务
func server(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Get response from %s server", addr)
	})
	s := http.Server{Addr: addr, Handler: mux}

	go func() {
		select {
		case <-ctx.Done():
			_ = s.Shutdown(ctx)
			fmt.Printf("server %s is cancelled\n", addr)
		}
	}()

	fmt.Printf("start %s service...\n", addr)
	return s.ListenAndServe()
}

// fakeService：模拟的服务处理函数
func fakeService(ctx context.Context) error {
	ch := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			close(ch)
			fmt.Println("fake service is cancelled")
		}
	}()

	fmt.Println("start fake service...")
	select {
	case <-time.After(10 * time.Second): //程序不中断的话默认10s后返回一个错误
		return errors.New("fake service error")
	case <-ch:
		return nil
	}
}

// listenSignal： 信号监听
func listenSignal(ctx context.Context) error {
	fmt.Println("start listen signal service...")
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh)

	go func() {
		select {
		case <-ctx.Done():
			signal.Stop(sigCh)
			close(sigCh)
			fmt.Println("listen signal service is cancelled")
		}
	}()

	if sig, ok := <-sigCh; ok {
		fmt.Printf("Got signal: %s\n", sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			return fmt.Errorf("Got %s signal, exit the program...", sig)
		default:
			return nil
		}
	} else {
		return nil
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return listenSignal(ctx)
	})
	g.Go(func() error {
		return server(ctx, "127.0.0.1:9000")
	})
	g.Go(func() error {
		return server(ctx, "127.0.0.1:9001")
	})
	g.Go(func() error {
		return fakeService(ctx)
	})
	if err := g.Wait(); err != nil {
		fmt.Printf("[main] Got error: %+v\n", err)
	}

	time.Sleep(2 * time.Second)
	fmt.Printf("[main] All services are cancelled!\n")
}
