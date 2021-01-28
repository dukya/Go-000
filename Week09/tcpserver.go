package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatal("listen error: %v", err)
	}

	var closing uint32
	atomic.StoreUint32(&closing, 0)

	exitSignals := []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, exitSignals...)
	go func() {
		for {
			s := <-sigCh
			log.Printf("Got signal: %s", s.String())
			switch s {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				atomic.StoreUint32(&closing, 1)
				listener.Close() // listener.Accept() will return err immediately
				return
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			if atomic.LoadUint32(&closing) != 0 {
				log.Printf("signal got, exit...")
				// 通知所有goroutine结束任务
				cancel()
				return
			}
			continue
		}

		//处理客户端连接
		go handleConn(ctx, conn)
	}
}

func handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	ch := make(chan string)

	ctx1, cancel := context.WithCancel(ctx)
	defer cancel()

	// 向客户端发送数据
	go sendMsg(ctx1, conn, ch)

	// 读取客户端数据
	input := bufio.NewScanner(conn)
	for {
		var msg string
		if input.Scan() {
			msg = input.Text()
		} else {
			if err := input.Err(); err != nil {
				log.Printf("read error: %v", err)
			}
			break
		}

		select {
		case <-ctx.Done():
			log.Printf("read ctx err: %+v", ctx.Err())
			return
		case ch <- msg:
		}
	}
}

//发送数据
func sendMsg(ctx context.Context, conn net.Conn, ch <-chan string) {
	wr := bufio.NewWriter(conn)
	for {
		select {
		case <-ctx.Done():
			log.Printf("writer ctx err: %+v", ctx.Err())
			return
		case msg := <-ch:
			log.Printf("send msg to client：%s", msg)
			wr.Write([]byte(msg))
			wr.Flush()
		}
	}
}
