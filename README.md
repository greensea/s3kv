# s3kv

A simple S3 compatible KV operation wrapper

[![Documentation](https://godoc.org/github.com/greensea/s3kv?status.svg)](http://godoc.org/github.com/greensea/s3kv)

## Usage

```go

// Init connection
s3, err = s3kv.New(&s3kv.Config{
	Endpoint:  os.Getenv("S3_ENDPOINT"),
	Bucket:    os.Getenv("S3_BUCKET"),
	AccessKey: os.Getenv("S3_ACCESS_KEY"),
	SecretKey: os.Getenv("S3_SECRET_KEY"),
})

// Put an object
s3.Put("foo", []byte("bar"))

// Get an object
bar, _ := s3.Get("foo")

// Marshal an object to JSON then Put
s3.PutObject("foo-json", []string{"foo", "bar"})

// Get an object and unmarshal it
foobar := make([]string)
s3.GetJSON("foo-json", &foobar)

// List key of objects
keys1, _ := s3.List("") // List all keys
keys2, _ := s3.List("prefix_") // List keys with "prefix_" prefix

// Delete an object
s3.Delete("foo")

// Check if a key is exists
isExists := s3.KeyExists("foo")
  
```
