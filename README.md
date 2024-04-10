# go-s3-tools
## usage
```
Usage of go-s3-tools:
  -accessKeyID string
        accessKeyID
  -bukcetName string
        bukcetName
  -endpoint string
        endpoint for storage
  -filePath string
        remove object name list
  -list_bucket
        list bucket object
  -need_prepare_data
        update 1000 object
  -prefix string
        list object, object prefix
  -secretAccessKey string
        secretAccessKey
  -test
        test if s3 is available
```

## now feature
### list bucket object
 ```
 go-s3-tools -list_bucket=true -endpoint=127.0.0.1:9000 -accessKeyID=minioadmin -secretAccessKey=minioadmin -bukcetName=test
 ```

### test the connectivity with S3
 ```
 go-s3-tools -test=true -endpoint=127.0.0.1:9000 -accessKeyID=minioadmin -secretAccessKey=minioadmin
 ```

### test program functional
The file in the `filePath` path does not need to contain data, it cann't be empty,
program will generate data and upload to s3, and write this object name to `filepath`
```
go-s3-tools -need_prepare_data=true -endpoint=127.0.0.1:9000 -accessKeyID=minioadmin -secretAccessKey=minioadmin -bukcetName=test -filePath=/root/code/go/go-s3-tools/file_list
```

### remove object from you provid lists
```
go-s3-tools -endpoint=127.0.0.1:9000 -accessKeyID=minioadmin -secretAccessKey=minioadmin -bukcetName=test -filePath=/root/code/go/go-s3-tools/file_list
```