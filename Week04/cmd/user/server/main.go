package main

import "fmt"

func main() {
	srv, cleanup, err := InitializeServer()
	defer cleanup()
	if err != nil {
		fmt.Printf("[InitServer] error:%v\n", err)
		return
	}

	fmt.Println("Start Server...")
	if err = srv.Run(); err != nil {
		fmt.Printf("[RunServer] error:%v\n", err)
		return
	}
}
