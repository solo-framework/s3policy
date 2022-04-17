package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"s3policy/internal/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Service struct {
	client *s3.S3
}

func NewService(conf *config.Config, profile string) *Service {

	settings := conf.GetSection(profile)

	awsCfg := aws.NewConfig()
	awsCfg.Region = aws.String(settings["region"])
	awsCfg.Endpoint = aws.String(settings["endpoint"])
	// awsCfg.HTTPClient.Timeout

	awsCfg.Credentials = credentials.NewStaticCredentials(settings["id"], settings["key"], "")

	sess := session.Must(session.NewSession(awsCfg))
	cl := s3.New(sess)

	return &Service{client: cl}
}

func (s *Service) GetPolicy(bucket string) error {

	pol, err := s.client.GetBucketPolicy(
		&s3.GetBucketPolicyInput{Bucket: aws.String(bucket)},
	)

	if err != nil {
		return fmt.Errorf("GetBucketPolicy error: %s", err)
		// panic(fmt.Sprintf("GetBucketPolicy error: %s", err))
	}

	fmt.Println(*pol.Policy)
	return nil
}

func (s *Service) DeletePolicy(bucket string) error {
	_, err := s.client.DeleteBucketPolicy(&s3.DeleteBucketPolicyInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return fmt.Errorf("unable to set bucket %q policy, %v", bucket, err)
	}

	fmt.Printf("Successfully deleted the policy on bucket %q.\n", bucket)
	return nil
}

func (s *Service) PutPolicy(bucket string, policyFile string) error {

	policyContent, err := ioutil.ReadFile(policyFile)
	if errors.Is(err, os.ErrNotExist) {
		// panic(fmt.Sprintf("Policy file doesn't exist: %s", err))
		return fmt.Errorf("policy file doesn't exist: %v", err)
	}

	var policyInput map[string]interface{}
	err = json.Unmarshal(policyContent, &policyInput)

	if err != nil {
		return fmt.Errorf("policy parsing error: %s", err)
	}

	_, err = s.client.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(string(policyContent)),
	})

	if err != nil {
		return fmt.Errorf("unable to set bucket %q policy, %v", bucket, err)
	}

	fmt.Printf("Successfully set bucket %q's policy\n", bucket)
	return nil
}

func (s *Service) CreateBucket(bucket string) error {

	_, err := s.client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return fmt.Errorf("unable to create bucket %q, %v", bucket, err)
	}

	fmt.Printf("Waiting for bucket %q to be created...\n", bucket)

	err = s.client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return fmt.Errorf("error occurred while waiting for bucket to be created, %v", bucket)
	}

	fmt.Printf("Bucket %q successfully created\n", bucket)
	return nil
}

func (s *Service) DeleteBucket(bucket string) error {

	_, err := s.client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return fmt.Errorf("unable to delete bucket %q, %v", bucket, err)
	}

	fmt.Printf("Waiting for bucket %q to be deleted...\n", bucket)

	err = s.client.WaitUntilBucketNotExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return fmt.Errorf("error occurred while waiting for bucket to be deleted, %v", bucket)
	}

	fmt.Printf("Bucket %q successfully deleted\n", bucket)
	return nil
}

func (s *Service) GetBucketACL(bucket string, verbose bool) error {

	result, err := s.client.GetBucketAcl(&s3.GetBucketAclInput{Bucket: &bucket})
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	if verbose {
		fmt.Println("Owner:", *result.Owner.ID)
		fmt.Println("")
		fmt.Println("ACL:")
	}
	fmt.Println(result)
	return nil
}

func (s *Service) ListBuckets() error {

	res, err := s.client.ListBuckets(
		&s3.ListBucketsInput{},
	)

	if err != nil {
		return fmt.Errorf("ListBuckets error: %s", err)
	}

	for _, v := range res.Buckets {
		fmt.Println(v)
	}
	return nil
}
