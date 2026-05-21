package routes

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go_project/models"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// BMComponent ===================== 【模型 动态结构体】 =====================
type BMComponent struct {
	ComponentID     string            `json:"component_id"`
	Name            string            `json:"name"`
	Manufacturer    string            `json:"manufacturer"`
	Size            string            `json:"size"`
	Description     string            `json:"description"`
	Author          string            `json:"author"`
	Version         string            `json:"version"`
	Created         string            `json:"created"`
	ModelParams     map[string]string `json:"model_params"`      // 动态参数
	ModelParamsUnit map[string]string `json:"model_params_unit"` // 动态单位
}

// UploadDir = "./upload"
var TempDir = "./temp"

func parseBMComponent(content string) BMComponent {
	var comp BMComponent
	comp.ModelParams = make(map[string]string)
	comp.ModelParamsUnit = make(map[string]string)

	// 1. 提取组件 UUID
	idRe := regexp.MustCompile(`\(bm_component\s+([a-f0-9-]+)`)
	if m := idRe.FindStringSubmatch(content); len(m) > 1 {
		comp.ComponentID = m[1]
	}

	// 2. 提取基础字段
	baseKv := map[string]*string{
		"name":         &comp.Name,
		"Manufacturer": &comp.Manufacturer,
		"Size":         &comp.Size,
		"description":  &comp.Description,
		"author":       &comp.Author,
		"version":      &comp.Version,
		"created":      &comp.Created,
	}
	for k, v := range baseKv {
		p := regexp.MustCompile(`\(` + regexp.QuoteMeta(k) + `\s+"([^"]+)"\)`)
		if m := p.FindStringSubmatch(content); len(m) > 1 {
			*v = m[1]
		}
	}

	// ===================== 【核心修复：全局提取所有 entry】 =====================
	// 匹配所有 (entry "key" "val")  —— 不管在哪个位置，全部提取
	entryRegex := regexp.MustCompile(`\(entry\s+"([^"]+)"\s+"([^"]+)"\)`)
	allEntries := entryRegex.FindAllStringSubmatch(content, -1)

	// 区分参数 / 单位（根据出现顺序）
	isParamSection := true
	for _, entry := range allEntries {
		key := entry[1]
		val := entry[2]

		// 第一个出现的是 ModelParams，后面的是 ModelParamsUnit
		if _, exists := comp.ModelParams[key]; exists {
			isParamSection = false
		}

		if isParamSection {
			comp.ModelParams[key] = val
		} else {
			comp.ModelParamsUnit[key] = val
		}
	}
	return comp
}

func uploadFile(c *gin.Context) {
	// 1. 获取上传的zip文件
	fileHeader, err := c.FormFile("zipfile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败：" + err.Error()})
		return
	}

	// 2. 打开文件流
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "打开文件失败"})
		return
	}
	defer file.Close()

	// 3. 读取到内存
	buf, _ := io.ReadAll(file)

	// 4. 解析ZIP
	zipReader, err := zip.NewReader(bytes.NewReader(buf), fileHeader.Size)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法ZIP文件"})
		return
	}

	// ===================== 核心：匹配动态目录 =====================
	// 固定前缀：filter_HPF/filter_HPF/cmp/
	basePath := "filter_HPF/cmp/"
	var epFiles []gin.H

	// 遍历压缩包内所有文件
	for _, f := range zipReader.File {
		filePath := f.Name
		// 跳过文件夹
		if f.FileInfo().IsDir() {
			continue
		}

		// 1. 必须在 basePath 路径下
		if !strings.HasPrefix(filePath, basePath) {
			continue
		}

		// 2. 必须是 .ep 后缀
		if filepath.Ext(filePath) != ".ep" {
			continue
		}

		// 3. 必须是 basePath 下一级目录里的文件（动态UUID目录）
		// 切割路径，确保结构是 basePath/xxx/xxx.ep
		relPath := strings.TrimPrefix(filePath, basePath)
		parts := strings.Split(relPath, "/")
		if len(parts) < 2 { // 必须在子目录里
			continue
		}

		// ===================== 读取 .ep 文件内容 =====================
		rc, err2 := f.Open()
		if err2 != nil {
			continue
		}
		content, _ := io.ReadAll(rc)
		rc.Close()

		// 收集结果
		epFiles = append(epFiles, gin.H{
			"file_path": filePath,        // 完整路径
			"uuid_dir":  parts[0],        // 动态UUID目录名
			"file_size": len(content),    // 文件大小
			"content":   string(content), // 文件内容
		})
		for i, e := range epFiles {
			fmt.Println(i)
			fmt.Println(e["file_path"])
			fmt.Println(e["uuid_dir"])
			fmt.Println(e["file_size"])
			comp := parseBMComponent(string(content))
			fmt.Println(comp)
			fmt.Println("模型ID：", comp.ComponentID)
			fmt.Println("描述：", comp.Description)
			fmt.Println("ModelParams ", comp.ModelParams["cutoff_freq"])
		}
	}

	// 保存文件到本地
	//savePath := "./upload/" + file.Filename
	//err = c.SaveUploadedFile(file, savePath)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error": "文件保存失败",
	//	})
	//	return
	//}
}

// 获取分片保存路径
func getChunkPath(fileMD5 string, index int) string {
	dir := filepath.Join(TempDir, fileMD5)
	os.MkdirAll(dir, 0755)
	return filepath.Join(dir, strconv.Itoa(index))
}

// 检查分片是否已存在
func isChunkExist(fileMD5 string, index int) bool {
	_, err := os.Stat(getChunkPath(fileMD5, index))
	return err == nil
}

// 检查所有分片是否上传完成
func isAllChunksUploaded(fileMD5 string, total int) bool {
	for i := 0; i < total; i++ {
		if !isChunkExist(fileMD5, i) {
			return false
		}
	}
	return true
}

func uploadChunkHandler(c *gin.Context) {
	fileMD5 := c.PostForm("file_md5")
	indexStr := c.PostForm("chunk_index")
	totalStr := c.PostForm("total_chunks")
	fmt.Println(fileMD5, indexStr, totalStr)
	if fileMD5 == "" || indexStr == "" || totalStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数不全"})
		return
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "chunk_index错误"})
		return
	}
	// 保存分片
	chunkFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "获取分片失败：" + err.Error(),
		})
		return
	}
	savePath := getChunkPath(fileMD5, index)
	if err2 := c.SaveUploadedFile(chunkFile, savePath); err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存分片失败"})
		return
	}
	var data models.UploadChunk
	data.CurrentChunk, _ = strconv.Atoi(indexStr)
	totalChunks, err3 := strconv.Atoi(totalStr)
	if err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "total_chunks错误"})
		return
	}
	if index == totalChunks-1 {
		data.IsAllUploaded = true
	} else {
		data.IsAllUploaded = false
	}
	msg := "分片" + indexStr + "上传成功"
	c.JSON(http.StatusCreated, gin.H{"code": 201, "msg": msg, "status": "success", "data": data})
}

func fileProcess(c *gin.Context) {

}
