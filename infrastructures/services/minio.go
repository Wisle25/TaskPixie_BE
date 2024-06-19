package services

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/wisle25/be-template/commons"
	"log"
)

// NewMinio make a new connection with Minio.
// Returning the client itself and the bucket name
func NewMinio(config *commons.Config) (*minio.Client, string) {
	ctx := context.Background()
	useSSL := config.AppEnv == "prod"
	var err error

	// Init
	minioClient, err := minio.New(config.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.MinioAccessKey, config.MinioSecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(fmt.Errorf("new minio client: init: %w", err))
	}

	// Make new bucket
	err = minioClient.MakeBucket(ctx, config.MinioBucket, minio.MakeBucketOptions{
		Region: config.MinioLocation,
	})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, config.MinioBucket)

		if errBucketExists == nil && exists {
			log.Printf("Bucket %s is already exists", config.MinioBucket)
		} else {
			panic(fmt.Errorf("create bucket %s: %w", config.MinioBucket, err))
		}
	} else {
		log.Printf("Bucket %s is created", config.MinioBucket)
	}

	return minioClient, config.MinioBucket
}
