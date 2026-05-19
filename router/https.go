package router

import (
	"go_project/controllers"
	"log"
	"path/filepath"
	"runtime"
)

func StartHTTPS() {
	log.Println("🚀 HTTPS服务启动: https://localhost:8002")
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	certFile := filepath.Join(currentDir, "server.pem")
	keyFile := filepath.Join(currentDir, "server.key")
	// 关键：在8002端口启动HTTPS
	err := controllers.R.RunTLS(":8002", certFile, keyFile)
	if err != nil {
		log.Fatal("❌ 启动失败:", err)
	}
}
