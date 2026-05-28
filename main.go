package main

import (
	"fmt"
	"go_project/router"
	"go_project/utils"
	"time"
)

//var ctx = context.Background()
//var rdb *redis.Client

func main() {
	fmt.Println("✅ 程序开始启动...")

	// 先捕获主函数的 panic，防止路由初始化崩溃
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("❌ 主程序崩溃: %v\n", err)
			// 让程序暂停一下，方便看错误
			time.Sleep(10 * time.Second)
		}
	}()

	fmt.Println("✅ 准备启动后台定时清理任务...")
	go utils.StartCleanShardJob()
	fmt.Println("✅ 后台定时清理任务已启动")

	fmt.Println("✅ 准备启动 HTTP 服务...")
	router.InitRouter()
	fmt.Println("❌ HTTP 服务意外退出了！")

	// 防止程序立刻退出，方便看错误
	time.Sleep(10 * time.Second)
}
