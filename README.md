## s3gopher
s3gopher is a light-weight S3 library using Go. It's built on top of [Amazon's offical AWS SDK](github.com/aws/aws-sdk-go/aws).

## Why
Found the AWS cumbersome to use if you just want to read / write a single or few files from a S3 bucket or list its content.

## Caveats 

Region is hard-coded to "eu-west-1" and ACL to "private". Probably not well suited for buckets with thousands or millions of files.

## Example usage

```go

import "github.com/datajet-io/s3gopher"

func main() {

bucket, err := s3gopher.New("myBucketName", "myAccessKey", "mySecretAccessKey")

// Are my credentials valid?

err := bucket.Test()

if err != nil {
	fmt.Printeln("Something bad happended when testing credentials.")
	return
}

// List content of the bucket 

fileList, err := o.List()

// Get the first file from the bucket

o, err := bucket.Get(fileList[0].Key)

if err != nil {
	fmt.Printeln("Something bad happended when reading the data.")
	return
}

// When was the file last modified?

fmt.Printeln(o.LastModified)

// Output its []byte content

fmt.Printeln(o.Data)

/*

Do something with the file's data
	
*/

// Write the file back to S3


err := bucket.Put(o)

if err != nil {
	fmt.Printeln("Something bad happended when writing the data.")
	return
}

// We're done

}


```





