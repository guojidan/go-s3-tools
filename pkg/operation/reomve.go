package operation

import (
	"go-s3-tools/pkg/client"
	"strings"
	"sync"
)

func Remove_object(ch <-chan string, wg *sync.WaitGroup, minioClient *client.Client, bucket_name *string) {
	defer wg.Done()

	for obj_name := range ch {
		trimed_obj_name := strings.TrimSpace(obj_name)
		// ignore error info
		go minioClient.Remove_object(bucket_name, &trimed_obj_name)
	}
}
