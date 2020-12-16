package grpctransport

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	v1 "geektime/Go-000/Week04/api/user/v1"
	"geektime/Go-000/Week04/internal/service"
)

type Server struct {
	service *service.UserService
}

func NewServer(svc *service.UserService) *Server {
	return &Server{service: svc}
}

func (srv *Server) Run() error {
	//TODO: 需要从配置文件读取
	listener, err := net.Listen("tcp", "127.0.0.1:9999")
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(context.Background())
	gs := grpc.NewServer()
	v1.RegisterUserServer(gs, srv.service)

	g.Go(func() error {
		go func() {
			<-ctx.Done()
			gs.GracefulStop()
			fmt.Println("Shutdown grpc server.")
		}()
		return gs.Serve(listener)
	})

	g.Go(func() error {
		exitSignals := []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
		sigCh := make(chan os.Signal, len(exitSignals))
		signal.Notify(sigCh, exitSignals...)
		for {
			select {
			case <-ctx.Done():
				return nil
			case s := <-sigCh:
				fmt.Printf("get a signal %s\n", s)
				return fmt.Errorf("Got %s signal, exit the program...\n", s)
			}
		}
	})

	return g.Wait()
}
