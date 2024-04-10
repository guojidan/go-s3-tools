package client

import (
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type Client struct {
	*minio.Client
}

func NewClient(endpoint *string, accessKeyID *string, secretAccessKey *string) (*Client, error) {
	// I think we don't need ssl
	useSSL := false

	var client Client
	// Initialize minio client object.
	minioClient, err := minio.New(*endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(*accessKeyID, *secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		logrus.Fatalln(err)
	}

	client.Client = minioClient

	return &client, err
}

func (client Client) Bucket_exist(bucket_name *string) (bool, error) {
	exist, err := client.BucketExists(context.Background(), *bucket_name)
	if err != nil {
		logrus.Fatalln(err)
		return false, err
	}
	if exist {
		logrus.Debugln("Bucket found")
	}

	return exist, err
}

func (client Client) Make_bucket(bucket_name *string) error {
	err := client.MakeBucket(context.Background(), *bucket_name, minio.MakeBucketOptions{Region: "", ObjectLocking: true})
	if err != nil {
		logrus.Fatalln(err)
		return err
	}

	return nil
}

func (client Client) Put_object(file_path *string, bucket_name *string, obj_name *string) error {
	file, err := os.Open(*file_path)
	if err != nil {
		logrus.Fatalln(err)
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		logrus.Fatalln(err)
		return err
	}

	uploadInfo, err := client.PutObject(context.Background(), *bucket_name, *obj_name, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		logrus.Fatalln(err)
		return err
	}
	logrus.Debugln("Successfully uploaded bytes: ", uploadInfo)
	return nil
}

func (client Client) List_object(bucket_name *string, prefix *string) error {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	objectCh := client.ListObjects(ctx, *bucket_name, minio.ListObjectsOptions{
		Prefix:    *prefix,
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			logrus.Fatalln(object.Err)
			return object.Err
		}
		logrus.Debugln(object.Key)
	}

	return nil
}

func (client Client) Remove_object(bucket_name *string, obj_name *string) error {
	opts := minio.RemoveObjectOptions{}
	err := client.RemoveObject(context.Background(), *bucket_name, *obj_name, opts)

	// ignore err, continue!
	if err != nil {
		logrus.Fatalln(err)
	}

	return nil
}
