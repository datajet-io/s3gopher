package s3gopher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
)

//Bucket handler for a S3 bucket
type Bucket struct {
	Client *s3.S3
	Bucket string
	ACL    string
	Creds  Credentials
}

// Credentials represents a pair of AWS credentials
type Credentials struct {
	AccessKey       string
	SecretAccessKey string
}

// Object in a S3 bucket
type Object struct {
	Key          string
	LastModified time.Time
	Data         []byte
}

func (p Object) String() string {
	return fmt.Sprintf("%s; %s", p.LastModified, p.Key)
}

// ByLastModified implements sort.Interface for []Object based on
// the LastModified field. Sort is descending by time
type ByLastModified []Object

func (a ByLastModified) Len() int           { return len(a) }
func (a ByLastModified) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLastModified) Less(i, j int) bool { return a[i].LastModified.After(a[j].LastModified) }

// New creates a new connection to a S3 bucket using the provided credentials.
func New(bucket, accessKey, secretAccessKey string) (*Bucket, error) {

	creds := credentials.NewStaticCredentials(accessKey, secretAccessKey, "")

	newS3 := s3.New(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: creds,
	})

	return &Bucket{Client: newS3, Bucket: bucket, ACL: "private"}, nil
}

//Test tests the connection to the bucket using the provided credentials
func (s *Bucket) Test() error {

	params := &s3.HeadBucketInput{
		Bucket: aws.String(s.Bucket), // Required
	}

	_, err := s.Client.HeadBucket(params)

	if err != nil {
		return err
	}

	return nil
}

//List all objects in a bucket
func (s *Bucket) List() (obj []Object, err error) {

	var resp *s3.ListObjectsOutput
	var currentKey string

	for {

		params := &s3.ListObjectsInput{
			Bucket:  aws.String(s.Bucket), // Required
			MaxKeys: aws.Int64(1000),
			Marker:  aws.String(currentKey),
		}
		resp, err = s.Client.ListObjects(params)

		for i := 0; i < len(resp.Contents); i++ {
			currentKey = *resp.Contents[i].Key
			obj = append(obj, Object{Key: *resp.Contents[i].Key, LastModified: *resp.Contents[i].LastModified})
		}

		if err != nil {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return nil, err
		}

		if !*resp.IsTruncated {
			break
		}
	}

	sort.Sort(ByLastModified(obj))

	return obj, nil
}

// Get data of an S3 object
func (s *Bucket) Get(key string) (o *Object, err error) {

	params := &s3.GetObjectInput{
		Bucket:          aws.String(s.Bucket),  // Required
		Key:             aws.String("/" + key), // Required
		ResponseExpires: aws.Time(time.Now()),
	}
	resp, err := s.Client.GetObject(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS Error with Code, Message, and original error (if any)
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				errMsg := fmt.Sprintf("%s %d %s", reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())

				return nil, errors.New(errMsg)
			}

			errMsg := fmt.Sprintf("%s %s %s", awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			return nil, errors.New(errMsg)

		}

		// This case should never be hit, the SDK should always return an
		// error which satisfies the awserr.Error interface.
		return nil, errors.New(err.Error())

	}

	buf, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return &Object{Key: key, LastModified: *resp.LastModified, Data: buf}, nil
}

//Put stores an object into an S3 bucket
func (s *Bucket) Put(o *Object) error {

	params := &s3.PutObjectInput{
		Bucket:      aws.String(s.Bucket),    // Required
		Key:         aws.String("/" + o.Key), // Required
		ACL:         aws.String(s.ACL),
		Body:        bytes.NewReader(o.Data),
		ContentType: aws.String("application/json"),
	}
	_, err := s.Client.PutObject(params)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Generic AWS Error with Code, Message, and original error (if any)
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				// A service error occurred
				errMsg := fmt.Sprintf("%s %d %s", reqErr.Message(), reqErr.StatusCode(), reqErr.RequestID())

				return errors.New(errMsg)
			}

			errMsg := fmt.Sprintf("%s %s %s", awsErr.Code(), awsErr.Message(), awsErr.OrigErr())
			return errors.New(errMsg)

		}

		// This case should never be hit, the SDK should always return an
		// error which satisfies the awserr.Error interface.
		return errors.New(err.Error())

	}

	return nil
}
