package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
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
			// we ignore big files
			continue
		}
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}

		randomId := utils.GenerateRandomId(16)
		mimeType := fileHeader.Header.Get("Content-Type")

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
			extension := "bin"
			if parts := strings.Split(mimeType, "/"); len(parts) == 2 {
				extension = parts[1]
			}
			key = fmt.Sprintf("attachment-%s.%s", randomId, extension)
		}

		as.uploadFile(key, mimeType, fileData, fileHeader.Filename)
		defer file.Close()

		attachmentUrl := fmt.Sprintf("%s/%s", as.cdnUrl, key)
		fileSize := utils.BytesToHuman(fileHeader.Size)

		attachment := Attachment{
			Id:       randomId,
			Url:      attachmentUrl,
			Filename: fileHeader.Filename,
			Filesize: fileSize,
			Type:     fileHeader.Header.Get("Content-Type"),
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
