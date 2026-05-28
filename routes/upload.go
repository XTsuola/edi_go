package routes

import (
	"archive/zip"
	"fmt"
	my "go_project/config"
	"go_project/models"
	"go_project/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

const chunkDir = "./temp"    // 分片存放目录
const mergeDir = "./uploads" // 合并后文件存放目录
// 解析模型压缩包
func processFile(c *gin.Context) {
	UUID := c.PostForm("file_uuid")
	if UUID == "" {
		c.JSON(400, gin.H{"error": "缺少 file_uuid"})
		return
	}

	// 1. 找到分片目录
	uuidChunkDir := filepath.Join(chunkDir, UUID)
	if _, err := os.Stat(uuidChunkDir); os.IsNotExist(err) {
		c.JSON(400, gin.H{"error": "分片目录不存在: " + uuidChunkDir})
		return
	}

	// 2. 读取分片
	entries, _ := os.ReadDir(uuidChunkDir)
	var chunks []string
	for _, e := range entries {
		if !e.IsDir() {
			chunks = append(chunks, filepath.Join(uuidChunkDir, e.Name()))
		}
	}

	// 3. 分片排序
	sort.Slice(chunks, func(i, j int) bool {
		a, _ := strconv.Atoi(filepath.Base(chunks[i]))
		b, _ := strconv.Atoi(filepath.Base(chunks[j]))
		return a < b
	})

	// 4. 合并成 ZIP
	mergePath := filepath.Join(mergeDir, UUID+".zip")
	dst, _ := os.Create(mergePath)
	for _, part := range chunks {
		src, _ := os.Open(part)
		io.Copy(dst, src)
		src.Close()
	}
	dst.Sync()
	dst.Close()

	// 5. 打开合并后的 ZIP
	zipFile, _ := os.Open(mergePath)
	fi, _ := zipFile.Stat()
	zipReader, _ := zip.NewReader(zipFile, fi.Size())
	defer zipFile.Close()

	type EpFile struct {
		Path    string `json:"path"`
		Content string `json:"content"` // 如果是二进制就用 []byte
	}

	// 你原来的变量
	var epFileList []EpFile

	// 新增：存储所有 uuid 目录 + 标记是否有 .ep 文件
	allUUIDDirs := make(map[string]bool)

	// --- 第一步：先把 cmp 下所有 uuid 文件夹找出来 ---
	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			parentDir := filepath.Dir(f.Name)
			if strings.HasSuffix(parentDir, "/cmp") || strings.HasSuffix(parentDir, "\\cmp") {
				allUUIDDirs[f.Name] = false // 先标记：没有ep
			}
		}
	}

	// --- 第二步：你原来的遍历逻辑（我只改一点点） ---
	for _, f := range zipReader.File {
		// 跳过文件夹
		if f.FileInfo().IsDir() {
			continue
		}

		// 只处理：.ep 后缀 + 在 cmp 目录下
		if strings.HasSuffix(f.Name, ".ep") && strings.Contains(f.Name, "/cmp/") {
			// 打开 .ep 文件
			rc, err := f.Open()
			if err != nil {
				continue
			}

			// 读取内容
			content, err := io.ReadAll(rc)
			rc.Close()

			if err == nil {
				epFileList = append(epFileList, EpFile{
					Path:    f.Name,
					Content: string(content),
				})
			}

			// --- 新增：标记这个 uuid 文件夹有 .ep ---
			fileDir := filepath.Dir(f.Name)
			if _, ok := allUUIDDirs[fileDir]; ok {
				allUUIDDirs[fileDir] = true
			}
		}
	}
	fmt.Println(allUUIDDirs, "888")

	// --- 第三步：找出没有 .ep 的 uuid ---
	var emptyUUIDs []string
	for dirPath, hasEp := range allUUIDDirs {
		if !hasEp {
			uuid := filepath.Base(dirPath) // 只拿 uuid 名称
			emptyUUIDs = append(emptyUUIDs, uuid)
		}
	}

	// --- 第四步：如果有空的，返回错误 ---
	var err error
	if len(emptyUUIDs) > 0 {
		err = fmt.Errorf("以下UUID文件夹缺少.ep文件：%v", emptyUUIDs)
	}

	// 查询categories_id最大值
	var model_packages_id int
	my.DB.Table("model_packages").Select("MAX(id)").Scan(&model_packages_id)
	nowTime := utils.NowTimestamptz()
	userID, _ := c.Get("userID")
	userIDStr, ok := userID.(string)
	if !ok {
		return
	}
	userUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(400, gin.H{"msg": "用户ID不是合法UUID"})
		return
	}
	data := models.ModelPackagesAll{
		ID:            model_packages_id + 1,
		FileUUID:      userIDStr + UUID, // 你接口里的 file_uuid
		FileName:      "blob",           // 文件名
		TotalChunks:   len(chunks),      // 总分片数
		FileMd5:       UUID,             // 后面可以计算
		StorageType:   "local",
		StoragePath:   filepath.Join(mergeDir, UUID+".zip"), // 合并后的文件路径
		FileSize:      0,                                    // 后面可以填真实大小
		Status:        "success",                            // 状态：已上传/已合并/已解析
		DeviceModelId: uuid.Nil,                             // 设备模型ID（如果有就传真实uuid）
		UploadedById:  userUUID,                             // 上传人ID
		CreatedTime:   nowTime,
		UpdatedTime:   nowTime,
	}
	if err2 := my.DB.Table("model_packages").Create(&data).Error; err2 != nil {
		MyErr(err2.Error(), c)
		return
	}
	CreateOk("新增成功", c)
	for _, epFile := range epFileList {
		comp := parseBMComponent(epFile.Content)
		fmt.Println(comp)
	}

	// 清理分片
	os.RemoveAll(uuidChunkDir)
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
	UUID := c.PostForm("file_uuid")
	indexStr := c.PostForm("chunk_index")
	totalStr := c.PostForm("total_chunks")
	if UUID == "" || indexStr == "" || totalStr == "" {
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
	savePath := getChunkPath(UUID, index)
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
