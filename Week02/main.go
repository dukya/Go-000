package main

import (
	"errors"
	"fmt"

	"geektime/Go-000/Week02/dao"
	"geektime/Go-000/Week02/service"
)

func main() {
	name := "test1"

	srv := &service.UserService{}
	res, err := srv.GetUser(name)

	//屏蔽掉底层的sql.ErrNoRows，使用自己的预定义错误
	if errors.Is(err, dao.ErrDaoNotFound) {
		fmt.Printf("%+v\n", err)
		//或者使用mock数据返回
		//mock := model.User{}
		//fmt.Printf("%+v\n", &mock)
		return
	}

	//其他错误则直接记录下来
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}

	//没有错误则打印结果
	fmt.Printf("%+v", res)
}
