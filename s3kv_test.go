package s3kv

import (
	"log"
	"os"
	"testing"
)

func TestS3kv(t *testing.T) {
	// Test for connection
	s, err := New(&Config{
		Endpoint:  os.Getenv("S3_ENDPOINT"),
		Bucket:    os.Getenv("S3_BUCKET"),
		AccessKey: os.Getenv("S3_ACCESS_KEY"),
		SecretKey: os.Getenv("S3_SECRET_KEY"),
	})
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	// Test for Put
	err = s.Put("s3kv_test_key", []byte(`"s3kv_test_value"`))
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	// Test for Get
	buf, err := s.Get("s3kv_test_key")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	// Test for GetJSON
	var j string
	err = s.GetJSON("s3kv_test_key", &j)
	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	// Test for PutObject
	err = s.PutObject("s3kv_test_key", map[string]interface{}{
		"Foo": "Bar",
	})

	if string(buf) != `"s3kv_test_value"` {
		t.FailNow()
	}

	// Test for KeyExists
	is, err := s.KeyExists("A_Not_Exists_Key___")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if is == true {
		log.Println("Key should not exists")
		t.FailNow()
	}

	is, err = s.KeyExists("s3kv_test_key")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if is != true {
		log.Println("Key should be exists")
		t.FailNow()
	}

	// Test for List
	res, err := s.List("s3kv_test_ke")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
	if len(res) < 1 {
		t.FailNow()
	}

	// Test for Delete
	err = s.Delete("s3kv_test_key")
	if err != nil {
		log.Println(err)
		t.FailNow()
	}
}
