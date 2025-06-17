package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/okzmo/kyob/internal/utils"
)

type AttachmentService struct {
	s3Client *s3.Client
	bucket   string
	cdnUrl   string
}

type Attachment struct {
	Id       string `json:"id"`
	Url      string `json:"url"`
	Filename string `json:"file_name"`
	Filesize string `json:"file_size"`
	Type     string `json:"type"`
}

func NewAttachmentService() *AttachmentService {
	return &AttachmentService{
		s3Client: s3.NewFromConfig(GetAWSConfig()),
		bucket:   "nyo-files",
		cdnUrl:   os.Getenv("CDN_URL"),
	}
}

func (as *AttachmentService) ProcessAttachments(files []*multipart.FileHeader, maxSize int64) ([]byte, error) {
	var attachments []Attachment

	for _, fileHeader := range files {
		if fileHeader.Size > maxSize {
			continue
		}

		if fileHeader.Size == 0 {
			continue
		}

		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}

		mimeType := http.DetectContentType(buffer[:n])

		if seeker, ok := file.(io.Seeker); ok {
			seeker.Seek(0, io.SeekStart)
		}

		randomId := utils.GenerateRandomId(16)
		var key string
		var fileData io.Reader = file

		if strings.Contains(mimeType, "image") {
			key = fmt.Sprintf("attachment-%s.webp", randomId)

			webpData, err := utils.ConvertToWebp(file)
			if err != nil {
				return nil, fmt.Errorf("failed to convert image to webp: %w", err)
			}
			fileData = bytes.NewReader(webpData)
		} else {
			extension := getSecureExtension(mimeType)
			key = fmt.Sprintf("attachment-%s.%s", randomId, extension)
		}

		if err := as.uploadFile(key, mimeType, fileData, fileHeader.Filename); err != nil {
			return nil, fmt.Errorf("failed to upload file: %w", err)
		}
		defer file.Close()

		attachmentUrl := fmt.Sprintf("%s/%s", as.cdnUrl, key)
		fileSize := utils.BytesToHuman(fileHeader.Size)

		attachment := Attachment{
			Id:       randomId,
			Url:      attachmentUrl,
			Filename: sanitizeFilename(fileHeader.Filename),
			Filesize: fileSize,
			Type:     mimeType,
		}

		attachments = append(attachments, attachment)
	}

	res, err := json.Marshal(attachments)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal attachments: %w", err)
	}

	return res, nil
}

func (as *AttachmentService) uploadFile(key string, mimeType string, fileData io.Reader, fileName string) error {
	input := &s3.PutObjectInput{
		Key:    &key,
		Bucket: aws.String("nyo-files"),
		Body:   fileData,
	}

	if !strings.Contains(mimeType, "image") && !strings.Contains(mimeType, "video") {
		input.ContentDisposition = aws.String(fmt.Sprintf(`attachment; filename="%s"`,
			strings.ReplaceAll(fileName, `"`, `\"`)))
	}

	_, err := as.s3Client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload attachment: %w", err)
	}

	return nil
}

func getSecureExtension(mimeType string) string {
	extensions := map[string]string{
		"application/pdf": "pdf",
		"text/plain":      "txt",
		"image/jpeg":      "jpg",
		"image/png":       "png",
	}

	if ext, ok := extensions[mimeType]; ok {
		return ext
	}

	if parts := strings.Split(mimeType, "/"); len(parts) == 2 {
		return parts[1]
	}

	return "bin"
}

func sanitizeFilename(filename string) string {
	filename = filepath.Base(filename)
	filename = strings.ReplaceAll(filename, "..", "")

	if len(filename) > 255 {
		filename = filename[:255]
	}

	return filename
}
