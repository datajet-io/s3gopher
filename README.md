## s3gopher
s3gopher is a light-weight S3 library using Go. It's built on top of [Amazon's offical AWS SDK](https://github.com/aws/aws-sdk-go/aws).

## Why
Found the official SDK to require a lot of boilerplate if you just want to read / write a single or few files from a S3 bucket or list its content.

## Caveats 

* See AWS documentation for [Region](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html) and [ACL](http://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#canned-acl) values.
* Listing is done in memory and not well suited for buckets with thousands or millions of files.

## Example usage

```go

import "github.com/datajet-io/s3gopher"

func main() {

config := s3gopher.Config{
		Bucket: "bucketname",
		ACL: "public-read",
		Region: "us-east-1",
		Credentials: s3gopher.Credentials{AccessKey: "yourAccessKey", SecretAccessKey:"yourSecretAccessKey"},
}

bucket, err := s3gopher.New(config)

// Are my credentials valid?

err := bucket.Test()

if err != nil {
	fmt.Printeln("Something bad happended when testing credentials.")
	return
}

// List content of the bucket, show latest additions first

fileList, err := o.List()

// Get the first file from the bucket

o, err := bucket.Get(fileList[0].Key)

if err != nil {
	fmt.Println("Something bad happended reading the data.")
	return
}

// When was the file last modified?

fmt.Println(o.LastModified)

// Output its []byte content

fmt.Println(o.Data)

/*

Do something with the file's data
	
*/

// Write the file back to S3


err := bucket.Put(o)

if err != nil {
	fmt.Println("Something bad happended when writing the data.")
	return
}

// We're done

}


```


## To Do

* Delete





