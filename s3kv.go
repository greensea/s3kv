package s3kv

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var ErrNotExists error = errors.New("Not exists")

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Region    string
	Bucket    string
}

type Storage struct {
	config  *Config
	session *session.Session
}

func New(config *Config) (*Storage, error) {
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, ""),
		Endpoint:    aws.String(config.Endpoint),
		Region:      aws.String("us-east-1"),
	})

	if err != nil {
		return nil, err
	}

	return &Storage{
		config:  config,
		session: sess,
	}, nil
}

func (s *Storage) Put(key string, value []byte) error {
	uploader := s3manager.NewUploader(s.session)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(value),
	})

	if err != nil {
		return err
	} else {
		return nil
	}
}

// Marshal an object and then put it to storage.
func (s *Storage) PutObject(key string, value any) error {
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.Put(key, buf)

}

// Get a key from storage.
// If key is not exists, a ErrNotExists error is returnd.
func (s *Storage) Get(key string) ([]byte, error) {
	downloader := s3manager.NewDownloader(s.session)

	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchKey {
				return nil, ErrNotExists
			}
		}
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}

// Get an JSON-encoded object and unmarshal it into value.
func (s *Storage) GetJSON(key string, value any) error {
	buf, err := s.Get(key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, value)

	return err
}

func (s *Storage) List(key_prefix string) ([]string, error) {
	svc := s3.New(s.session)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(s.config.Bucket),
		Prefix: aws.String(key_prefix),
	}
	resp, err := svc.ListObjects(params)
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(resp.Contents))
	for k, item := range resp.Contents {
		ret[k] = string(*item.Key)
	}

	return ret, nil
}

func (s *Storage) Delete(key string) error {
	svc := s3.New(s.session)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(s.config.Bucket), Key: aws.String(key)})
	if err != nil {
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})

	return nil
}

// Check if a key is exists.
// The first return value indicate if the key is exists.
// If failed, an error is returned and the first return value is meaningless.
func (s *Storage) KeyExists(key string) (bool, error) {
	_, err := s.Get(key)
	if err == ErrNotExists {
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		return false, err
	}
}
