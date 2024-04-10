package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"go-s3-tools/pkg/client"
	"go-s3-tools/pkg/file"
	"go-s3-tools/pkg/operation"
)

func test_func(minioClient *client.Client) {
	file_path := "/root/code/go/go-s3-tools/go.mod"
	bucket_name := "testbucket"
	obj_name := "testobj"

	// check if the bucket exists
	exist, err := minioClient.Bucket_exist(&bucket_name)
	if err != nil {
		return
	}
	if !exist {
		// should make bucket
		err = minioClient.Make_bucket(&bucket_name)
		if err != nil {
			return
		}
	}

	// put object
	err = minioClient.Put_object(&file_path, &bucket_name, &obj_name)
	if err != nil {
		return
	}

	// remove object
	err = minioClient.Remove_object(&bucket_name, &obj_name)
	if err != nil {
		return
	}
}

func run(minioClient *client.Client, bucket_name *string, file_path *string) {
	// check if the bucket exists
	exist, err := minioClient.Bucket_exist(bucket_name)
	if err != nil {
		return
	}

	if !exist {
		log.Fatalln("bucket don't exist, please check")
		return
	}

	// I think we need a buffered channel
	ch := make(chan string, 4096)
	var wg sync.WaitGroup

	wg.Add(1)
	go file.Read_line(ch, &wg, file_path)

	num_thread := 3
	wg.Add(num_thread)
	for num_thread > 0 {
		go operation.Remove_object(ch, &wg, minioClient, bucket_name)
		num_thread--
	}
	wg.Wait()
}

func main() {
	// parse args

	endpoint := flag.String("endpoint", "", "endpoint for storage")
	accessKeyID := flag.String("accessKeyID", "", "accessKeyID")
	secretAccessKey := flag.String("secretAccessKey", "", "secretAccessKey")
	bukcetName := flag.String("bukcetName", "", "bukcetName")
	filePath := flag.String("filePath", "", "remove object name list")
	test := flag.Bool("test", false, "test if s3 is available")
	need_prepare_data := flag.Bool("need_prepare_data", false, "update 1000 object")
	list_bucket := flag.Bool("list_bucket", false, "list bucket object")
	prefix := flag.String("prefix", "", "list object, object prefix")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of go-s3-tools:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *endpoint == "" || *accessKeyID == "" || *secretAccessKey == "" {
		log.Fatalln("args err, please run '-h' check usage")
		return
	}

	minioClient, err := client.NewClient(endpoint, accessKeyID, secretAccessKey)
	if err != nil {
		return
	}

	if *test {
		test_func(minioClient)
		return
	}

	if *list_bucket {
		minioClient.List_object(bukcetName, prefix)
		return
	}

	if *need_prepare_data {
		file.Prepare_test_data(minioClient, bukcetName, filePath)
	}

	run(minioClient, bukcetName, filePath)

	log.Println("finish op")
}
