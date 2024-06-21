package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"

	"go-s3-tools/pkg/client"
	"go-s3-tools/pkg/file"
	"go-s3-tools/pkg/operation"

	"github.com/sirupsen/logrus"
)

// config
type Config struct {
	EndPoint 		string 	`json:"EndPoint"`
	AccessKeyID 	string 	`json:"AccessKeyID"`
	SecretAccessKey string 	`json:"SecretAccessKey"`
	BukcetName 		string 	`json:"BukcetName"`
	FilePath 		string 	`json:"FilePath"`
	Debug 			bool 	`json:"Debug"`
	Test 			bool 	`json:"Test"`
	NeedPrepareData bool 	`json:"NeedPrepareData"`
	ListBucket 		bool 	`json:"ListBucket"`
	Prefix 			string 	`json:"Prefix"`
	NumThread 		int 	`json:"NumThread"`
}

// global config
var config Config

func test_func(minioClient *client.Client) {
	file_path := "./test_file"
	bucket_name := "testbucket"
	obj_name := "testobj"

	err := file.Write_data(&file_path)
	if err != nil {
		logrus.Fatalln("write test data failed: {}", err)
		return
	}

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

func run(minioClient *client.Client, bucket_name *string, file_path *string, num_thread int) {
	// check if the bucket exists
	exist, err := minioClient.Bucket_exist(bucket_name)
	if err != nil {
		return
	}

	if !exist {
		logrus.Fatalln("bucket don't exist, please check")
		return
	}

	// I think we need a buffered channel
	ch := make(chan string, 4096)
	var wg sync.WaitGroup

	wg.Add(1)
	go file.Read_line(ch, &wg, file_path)

	wg.Add(num_thread)
	for num_thread > 0 {
		go operation.Remove_object(ch, &wg, minioClient, bucket_name)
		num_thread--
	}
	wg.Wait()
}

func load_config_from_file(config_file *string) error {
	file, err := os.Open(*config_file)
    if err != nil {
        fmt.Println("open config file error please check, Error:", err)
        return err
    }
    defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Println("can not decode config file, Err: ", err)
		return err
	}

	return nil
}

func main() {
	// parse args
	config_file := flag.String("config_file", "", "config file path")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of go-s3-tools:\n")
		flag.PrintDefaults()
	}
	
	flag.Parse()
	
	if *config_file != "" {
		err := load_config_from_file(config_file)
		if err != nil {
			return
		}
	}

	if config.EndPoint == "" || config.AccessKeyID == "" || config.SecretAccessKey == "" {
		logrus.Fatalln("args err, please run '-h' check usage")
		return
	}

	if config.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.FatalLevel)
	}

	minioClient, err := client.NewClient(&config.EndPoint, &config.AccessKeyID, &config.SecretAccessKey)
	if err != nil {
		return
	}

	if config.Test {
		test_func(minioClient)
		return
	}

	if config.ListBucket {
		minioClient.List_object(&config.BukcetName, &config.Prefix)
		return
	}

	if config.NeedPrepareData {
		file.Prepare_test_data(minioClient, &config.BukcetName, &config.FilePath)
	}

	run(minioClient, &config.BukcetName, &config.FilePath, config.NumThread)

	logrus.Infoln("finish op")
}
