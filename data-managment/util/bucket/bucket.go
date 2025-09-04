package bucket

import (
	"context"
	"data-managment/util/env"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var Bucket *MinioBucket

const bucketName = "scraped-data"

type MinioBucket struct {
	minioClientInstance *minio.Client
	lock                sync.Mutex
}

func Setup() error {
	zap.S().Debugf("Setting up MinIO storage")
	minioClient, err := minio.New(env.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(env.AccessKeyID, env.SecretAccessKey, ""),
		Secure: env.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("failed to create MinIO client %w", err)
	}

	Bucket = &MinioBucket{
		minioClientInstance: minioClient,
	}

	Bucket.lock.Lock()
	defer Bucket.lock.Unlock()
	ok, err := Bucket.minioClientInstance.BucketExists(context.Background(), bucketName)
	if err != nil {
		zap.S().Errorf("error pinging minio service with config(%s, %s)\n%s\n", env.Endpoint, env.AccessKeyID, err.Error())
		return err
	}

	if !ok {
		zap.S().Infof("Making new bucket: %s", bucketName)
		err = Bucket.
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

func (b *MinioBucket) Upload(file multipart.File, header *multipart.FileHeader) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	_, err := b.minioClientInstance.PutObject(
		context.Background(),
		bucketName,
		header.Filename,
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return fmt.Errorf("failed to upload %w", err)
	}

	return nil
}

func (b *MinioBucket) GetAllReaders() ([]io.ReadCloser, error) {
	b.lock.Lock()
	defer b.lock.Unlock()
	readers := make([]io.ReadCloser, 0)

	for file := range b.minioClientInstance.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{}) {
		zap.S().Debugf("Retriving object with name %s", file.Key)
		reader, err := b.minioClientInstance.GetObject(context.Background(), bucketName, file.Key, minio.GetObjectOptions{})
		if err != nil {
			return nil, err
		}
		readers = append(readers, reader)
	}
	return readers, nil
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

func (b *MinioBucket) UploadMany(files []*os.File) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	zap.S().Debugf("Starting file upload, files %d", len(files))

	count := 0
	for _, file := range files {
		info, _ := file.Stat()
		zap.S().Debugf("Uploading file name %s", info.Name())
		filename := filepath.Base(file.Name())

		_, err := b.minioClientInstance.PutObject(
			context.Background(),
			bucketName,
			filename,
			file,
			info.Size(),
			minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if err != nil {
			zap.S().Errorf("Failed to upload err = %v", err)
			continue
		}
		count++
	}

	return count, nil
}
