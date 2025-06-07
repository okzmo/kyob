package services

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/okzmo/kyob/internal/utils"
)

type AttachmentService struct {
	s3Client *s3.Client
	bucket   string
	cdnUrl   string
}

func NewAttachmentService() *AttachmentService {
	return &AttachmentService{
		s3Client: s3.NewFromConfig(GetAWSConfig()),
		bucket:   "nyo-files",
		cdnUrl:   os.Getenv("CDN_URL"),
	}
}

func (as *AttachmentService) ProcessAttachments(files []*multipart.FileHeader, maxSize int64) ([]string, error) {
	var attachmentsURLs []string

	for _, fileHeader := range files {
		if fileHeader.Size > maxSize {
			return nil, fmt.Errorf("file %s exceeds maximum size", fileHeader.Filename)
		}
		file, err := fileHeader.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %w", err)
		}

		randomId := utils.GenerateRandomId(16)
		fileName := fmt.Sprintf("attachment-%s.webp", randomId)

		_, err = as.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Key:    &fileName,
			Bucket: aws.String("nyo-files"),
			Body:   file,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to upload attachment: %w", err)
		}
		defer file.Close()

		attachmentUrl := fmt.Sprintf("%s/%s", as.cdnUrl, fileName)
		attachmentsURLs = append(attachmentsURLs, attachmentUrl)
	}

	return attachmentsURLs, nil
}
