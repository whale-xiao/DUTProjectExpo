package service

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"showcase-backend/pkg/errcode"
)

// UploadService 图片上传（V1 本地磁盘；二期换 OSS 时替换保存实现即可）。
type UploadService struct {
	dir string // UploadDir，如 ./uploads
}

func NewUploadService(dir string) *UploadService { return &UploadService{dir: dir} }

const maxUploadSize = 5 << 20 // 5MB

var allowedExt = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}

// Save 保存上传文件到 <UploadDir>/<yyyyMM>/<rand>.<ext>，返回可访问相对路径。
func (s *UploadService) Save(file *multipart.FileHeader) (string, error) {
	if file.Size > maxUploadSize {
		return "", NewError(errcode.InvalidParam, "图片不能超过 5MB")
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExt[ext] {
		return "", NewError(errcode.InvalidParam, "仅支持 jpg/png/webp/gif 图片")
	}

	sub := time.Now().Format("200601")
	dir := filepath.Join(s.dir, sub)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", NewError(errcode.InternalError, "存储目录创建失败")
	}

	name := randHex(16) + ext
	dst := filepath.Join(dir, name)

	src, err := file.Open()
	if err != nil {
		return "", NewError(errcode.InternalError, "读取上传文件失败")
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", NewError(errcode.InternalError, "保存文件失败")
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return "", NewError(errcode.InternalError, "保存文件失败")
	}
	return "/uploads/" + sub + "/" + name, nil
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// 几乎不可能失败；失败则退回时间戳，避免接口不可用
		return time.Now().Format("150405.000000000")
	}
	return hex.EncodeToString(b)
}
