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
	GetFiles() ([]io.ReadCloser, error)
	CheckBucket(name string) bool
	UploadMany(files []*multipart.FileHeader) (int, error)
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

func (b *MinioBucket) GetFiles() ([]io.ReadCloser, error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	ctx := context.Background()
	readers := make([]io.ReadCloser, 0)
	objectCh := b.minioClientInstance.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})

	for file := range objectCh {
		if file.Err != nil {
			zap.S().Errorf("Failed to list object: %v", file.Err)
			for _, r := range readers {
				_ = r.Close()
			}
			return nil, file.Err
		}

		zap.S().Debugf("Retrieving object with name %s", file.Key)
		reader, err := b.minioClientInstance.GetObject(ctx, bucketName, file.Key, minio.GetObjectOptions{})
		if err != nil {
			zap.S().Errorf("Failed to get object '%s': %v", file.Key, err)
			for _, r := range readers {
				_ = r.Close()
			}
			return nil, err
		}
		readers = append(readers, reader)
	}

	return readers, nil
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

func (b *MinioBucket) UploadMany(files []*multipart.FileHeader) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	zap.S().Debugf("Starting file upload, files %d", len(files))

	count := 0
	for _, file := range files {
		fileReader, err := file.Open()
		if err != nil {
			zap.S().Errorf("Failed to open file header err = %v", err)
			continue
		}
		defer fileReader.Close()

		zap.S().Debugf("Uploading file name %s", file.Filename)
		filename := filepath.Base(file.Filename)

		_, err = b.minioClientInstance.PutObject(
			context.Background(),
			bucketName,
			filename,
			fileReader,
			file.Size,
			minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
		if err != nil {
			zap.S().Errorf("Failed to upload err = %v", err)
			continue
		}
		count++
	}

	return count, nil
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
