package s3

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"time"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.peniremit.dev/peniremit-shared-upload-service/internal/port"
)

type S3Provider struct {
	client *s3.Client
	bucket string
}

func New(client *s3.Client, bucket string) *S3Provider {
	return &S3Provider{
		client: client,
		bucket: bucket,
	}
}

func (p *S3Provider) PresignUpload(ctx context.Context, in port.UploadInput, expires time.Duration) (port.PresignResponse, error) {
	presigner := s3.NewPresignClient(p.client)

	putInput := &s3.PutObjectInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(in.Key),
		ContentType: aws.String(in.ContentType),
	}

	if len(in.Metadata) > 0 {
		putInput.Metadata = convertMap(in.Metadata)
	}

	res, err := presigner.PresignPutObject(ctx, putInput, s3.WithPresignExpires(expires))
	if err != nil {
		return port.PresignResponse{}, fmt.Errorf("s3 presign put: %w", err)
	}

	return port.PresignResponse{
		URL:     res.URL,
		Fields:  nil,
		Expires: expires,
		Key:     in.Key,
	}, nil
}

func (p *S3Provider) PresignDownload(ctx context.Context, key string, expires time.Duration) (port.PresignResponse, error) {
	presigner := s3.NewPresignClient(p.client)

	out, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expires))
	if err != nil {
		return port.PresignResponse{}, fmt.Errorf("s3 presign get: %w", err)
	}

	return port.PresignResponse{
		URL:     out.URL,
		Expires: expires,
		Key:     key,
	}, nil
}

func (p *S3Provider) DeleteObject(ctx context.Context, key string) error {
	_, err := p.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(p.bucket),
		Key:    aws.String(key),
	})
	return err
}

func (p *S3Provider) PresignMultiUpload(ctx context.Context, inputs []port.UploadInput, expires time.Duration) ([]port.PresignResponse, error) {
	out := make([]port.PresignResponse, 0, len(inputs))
	for _, in := range inputs {
		pr, err := p.PresignUpload(ctx, in, expires)
		if err != nil {
			return nil, err
		}
		out = append(out, pr)
	}
	return out, nil
}

func (p *S3Provider) InitiateMultipart(ctx context.Context, in port.UploadInput) (port.MultipartSession, error) {
	res, err := p.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(p.bucket),
		Key:         aws.String(in.Key),
		ContentType: aws.String(in.ContentType),
		Metadata:    convertMap(in.Metadata),
	})
	if err != nil {
		return port.MultipartSession{}, fmt.Errorf("create multipart upload: %w", err)
	}

	return port.MultipartSession{
		UploadID: aws.ToString(res.UploadId),
		Key:      in.Key,
		Expires:  time.Now().Add(24 * time.Hour),
	}, nil
}

func (p *S3Provider) PresignMultipartPart(ctx context.Context, key, uploadID string, partNumber int32, expires time.Duration) (port.PresignResponse, error) {
	presigner := s3.NewPresignClient(p.client)

	in := &s3.UploadPartInput{
		Bucket:     aws.String(p.bucket),
		Key:        aws.String(key),
		UploadId:   aws.String(uploadID),
		PartNumber: aws.Int32(partNumber),
	}

	res, err := presigner.PresignUploadPart(ctx, in, s3.WithPresignExpires(expires))
	if err != nil {
		return port.PresignResponse{}, fmt.Errorf("presign upload part: %w", err)
	}

	return port.PresignResponse{
		URL:     res.URL,
		Expires: expires,
		Key:     key,
	}, nil
}

func (p *S3Provider) CompleteMultipart(ctx context.Context, key, uploadID string, parts []port.UploadPart) error {
	if uploadID == "" {
		return errors.New("uploadID required")
	}
	if len(parts) == 0 {
		return errors.New("no parts provided")
	}

	var completedParts []s3types.CompletedPart
	for _, pt := range parts {
		completedParts = append(completedParts, s3types.CompletedPart{
			ETag:       aws.String(pt.ETag),
			PartNumber: aws.Int32(pt.PartNumber),
		})
	}

	_, err := p.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(p.bucket),
		Key:             aws.String(key),
		UploadId:        aws.String(uploadID),
		MultipartUpload: &s3types.CompletedMultipartUpload{Parts: completedParts},
	})

	return err
}

// AbortMultipart aborts a multipart upload.
func (p *S3Provider) AbortMultipart(ctx context.Context, key, uploadID string) error {
	if uploadID == "" {
		return errors.New("uploadID required")
	}
	_, err := p.client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(p.bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
	})
	return err
}

func convertMap(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	out := make(map[string]string, len(src))
	maps.Copy(out, src)
	return out
}