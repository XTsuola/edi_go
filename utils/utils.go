package utils

import (
	"encoding/base64"
	"encoding/json"
	"os"
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
