package service

import (
	"PatientManager/config"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type IbucketService interface {
	CheckBucket(name string) bool
	UploadMany(files []*multipart.FileHeader, namePrefix string) ([]string, error)
	GetFile(name string) (io.ReadCloser, error)
	DeleteMany(names []string) error
}

var bucket *MinioBucket

const bucketName = "checkup-images"

type MinioBucket struct {
	minioClientInstance *minio.Client
	lock                sync.Mutex
}

func NewBucketService() IbucketService {
	if bucket == nil {
		setup()
	}
	return bucket
}

func setup() error {
	zap.S().Debugf("Setting up MinIO storage")
	minioClient, err := minio.New(config.AppConfig.MIOEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AppConfig.MIOAccessKeyID, config.AppConfig.MIOSecretAccessKey, ""),
		Secure: config.AppConfig.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client %w", err)
	}

	bucket = &MinioBucket{
		minioClientInstance: minioClient,
	}

	bucket.lock.Lock()
	defer bucket.lock.Unlock()
	ok, err := bucket.minioClientInstance.BucketExists(context.Background(), bucketName)
	if err != nil {
		zap.S().Errorf("error pinging minio service with config(%s, %s)\n%s\n", config.AppConfig.MIOEndpoint, config.AppConfig.MIOAccessKeyID, err.Error())
		return err
	}

	if !ok {
		zap.S().Infof("Making new bucket: %s", bucketName)
		err = bucket.
			minioClientInstance.
			MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			zap.S().Errorf("error making new bucket: %s, error: %v", bucketName, err)
			return err
		}
	}

	zap.S().Debugf("MinIO storage setup")
	return nil
}

func (b *MinioBucket) GetFile(name string) (io.ReadCloser, error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	ctx := context.Background()

	zap.S().Debugf("Retrieving object with name %s", name)
	reader, err := b.minioClientInstance.GetObject(ctx, bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		zap.S().Errorf("Failed to get object '%s': %v", name, err)
		return nil, err
	}
	return reader, nil
}

func (b *MinioBucket) CheckBucket(name string) bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	ret, err := b.minioClientInstance.BucketExists(context.Background(), name)
	if err != nil {
		zap.S().Errorf("Error accessing bucket: %s\n%s\n", name, err.Error())
		return false
	}
	return ret
}

func (b *MinioBucket) UploadMany(files []*multipart.FileHeader, namePrefix string) ([]string, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	zap.S().Debugf("Starting file upload for prefix %s, files %d", namePrefix, len(files))

	var uploadedPaths []string
	var uploadErrors []string

	for _, file := range files {
		fileReader, err := file.Open()
		if err != nil {
			msg := fmt.Sprintf("failed to open file header for %s: %v", file.Filename, err)
			zap.S().Error(msg)
			uploadErrors = append(uploadErrors, msg)
			continue
		}
		defer fileReader.Close()

		originalFilename := filepath.Base(file.Filename)
		newFilename := fmt.Sprintf("%s_%s", namePrefix, originalFilename)
		zap.S().Debugf("Uploading file with new name: %s", newFilename)

		_, err = b.minioClientInstance.PutObject(
			context.Background(),
			bucketName,
			newFilename,
			fileReader,
			file.Size,
			minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")},
		)
		if err != nil {
			msg := fmt.Sprintf("failed to upload %s: %v", newFilename, err)
			zap.S().Error(msg)
			uploadErrors = append(uploadErrors, msg)
			continue
		}
		uploadedPaths = append(uploadedPaths, newFilename)
	}

	if len(uploadErrors) > 0 {
		return uploadedPaths, fmt.Errorf("encountered errors during upload: %v", uploadErrors)
	}

	return uploadedPaths, nil
}

func (b *MinioBucket) DeleteMany(names []string) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if len(names) == 0 {
		return nil
	}

	zap.S().Debugf("Attempting to delete %d objects from bucket '%s'", len(names), bucketName)

	objectsCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectsCh)
		for _, name := range names {
			objectsCh <- minio.ObjectInfo{Key: name}
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	errorCh := b.minioClientInstance.RemoveObjects(context.Background(), bucketName, objectsCh, opts)

	var deleteErrors []string
	for e := range errorCh {
		if e.Err != nil {
			errMsg := fmt.Sprintf("Failed to remove object '%s', error: %v", e.ObjectName, e.Err)
			zap.S().Error(errMsg)
			deleteErrors = append(deleteErrors, errMsg)
		}
	}

	if len(deleteErrors) > 0 {
		return fmt.Errorf("encountered errors during object deletion: %v", deleteErrors)
	}

	zap.S().Infof("Successfully deleted %d objects from bucket '%s'", len(names), bucketName)
	return nil
}
