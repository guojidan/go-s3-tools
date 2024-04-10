package operation

import (
	"go-s3-tools/pkg/client"
	"log"
	"strings"
	"sync"
)

func Remove_object(ch <-chan string, wg *sync.WaitGroup, minioClient *client.Client, bucket_name *string) {
	defer wg.Done()

	for obj_name := range ch {
		trimed_obj_name := strings.TrimSpace(obj_name)
		err := minioClient.Remove_object(bucket_name, &trimed_obj_name)
		if err != nil {
			log.Fatalln("delete object faile: {}, err: {}", obj_name, err)
		}
	}
}