package main

import (
	"context"
	"fmt"
	v1 "geektime/Go-000/Week04/api/user/v1"
	"google.golang.org/grpc"
)

func main() {
	//1、连接rpc服务器
	conn, err := grpc.Dial("127.0.0.1:9999", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	//2、创建客户端并调用相应方法获取返回值
	client := v1.NewUserClient(conn)
	res, err := client.GetUserById(context.Background(), &v1.GetUserByIdRequest{Id: 1})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
