package main

import (
	"hotNews/db"
	"hotNews/routes"
)

func main() {
	defer db.DbClose()

	//channel
	//channel.Init()

	// 初始化路由
	routes.Init()
}
