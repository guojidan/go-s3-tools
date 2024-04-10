package file

import (
	"bufio"
	"fmt"
	"go-s3-tools/pkg/client"
	"log"
	"os"
)

func Prepare_test_data(minioClient *client.Client, bucket_name *string, file_path *string) error {
	// check if the bucket exists
	exist, err := minioClient.Bucket_exist(bucket_name)
	if err != nil {
		return err
	}
	if !exist {
		// should make bucket
		err = minioClient.Make_bucket(bucket_name)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(*file_path)
	if err != nil {
		log.Fatalln("create file err {}, err: {}", file_path, err)
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	object_base_name := "test_object"
	num_object := 1000
	for x := range num_object {
		object_name := object_base_name + fmt.Sprintf("%d", x)
		_, err = writer.WriteString(object_name + "\n")
		if err != nil {
			fmt.Println("cann't write data to file: ", err)
			return err
		}
		minioClient.Put_object(file_path, bucket_name, &object_name)
	}
	err = writer.Flush()
	if err != nil {
		fmt.Println("cann't flush buffer:", err)
		return err
	}

	return nil
}
