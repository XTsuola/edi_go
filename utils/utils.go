package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// WriteImg 写入二进制图片文件封装
func WriteImg(imgStr string, savePath string) {
	if idx := strings.Index(imgStr, ","); idx != -1 {
		imgStr = imgStr[idx+1:]
	}
	data, _ := base64.StdEncoding.DecodeString(imgStr)
	err := os.WriteFile(savePath, data, 0644)
	if err != nil {
		return
	}
}

// NowTimestamptz 获取当前时间
func NowTimestamptz() pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  time.Now(),
		Valid: true,
	}
}

// If 三元表达式
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// ArrToString 数组转字符串
func ArrToString[T any](arr []T) string {
	if len(arr) == 0 {
		return `[]`
	} else {
		jsonBytes, _ := json.Marshal(arr)
		jsonStr := string(jsonBytes)
		return jsonStr
	}
}

// StringToArr 字符串转数组
func StringToArr[T any](str string) []T {
	var arr []T
	err := json.Unmarshal([]byte(str), &arr)
	if err != nil || len(arr) == 0 {
		arr = []T{}
	}
	return arr
}

// StartCleanShardJob 核心：定时清理废弃分片（每8小时执行一次）
func StartCleanShardJob() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("❌ 定时任务崩溃: %v\n", err)
			// 崩溃后 10 秒自动重启定时任务
			fmt.Println("🔄 10秒后重启定时任务...")
			time.Sleep(10 * time.Second)
			go StartCleanShardJob()
		}
	}()

	fmt.Println("✅ 定时任务已进入运行状态")
	ticker := time.NewTicker(8 * time.Hour)
	defer ticker.Stop()

	fmt.Println("✅ 执行第一次分片清理...")
	CleanExpiredShards()
	fmt.Println("✅ 第一次分片清理完成")

	for range ticker.C {
		fmt.Println("⏰ 执行定时分片清理...")
		CleanExpiredShards()
		fmt.Println("✅ 定时分片清理完成")
	}
}

// CleanExpiredShards 真正执行清理的逻辑
func CleanExpiredShards() {
	shardDir := "./temp"
	expireDuration := 24 * time.Second

	// 自动创建 temp 目录（如果不存在）
	if _, err := os.Stat(shardDir); os.IsNotExist(err) {
		fmt.Printf("📁 分片目录 %s 不存在，自动创建\n", shardDir)
		err2 := os.MkdirAll(shardDir, 0755)
		if err2 != nil {
			fmt.Printf("❌ 创建分片目录失败: %v\n", err2)
			return
		}
	}

	entries, err := os.ReadDir(shardDir)
	if err != nil {
		fmt.Printf("❌ 读取分片目录失败: %v\n", err)
		return
	}

	fmt.Printf("📋 找到 %d 个分片文件/目录\n", len(entries))

	now := time.Now()
	deletedCount := 0

	for _, entry := range entries {
		info, err2 := entry.Info()
		if err2 != nil {
			fmt.Printf("❌ 获取文件信息失败: %s, %v\n", entry.Name(), err2)
			continue
		}

		if now.Sub(info.ModTime()) > expireDuration {
			delPath := filepath.Join(shardDir, entry.Name())
			err := os.RemoveAll(delPath)
			if err != nil {
				fmt.Printf("❌ 删除失败: %s, %v\n", delPath, err)
			} else {
				fmt.Printf("🗑️  已删除过期分片: %s\n", delPath)
				deletedCount++
			}
		}
	}

	fmt.Printf("✅ 本次清理完成，共删除 %d 个过期分片\n", deletedCount)
}
