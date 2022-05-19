package main

import (
	"fmt"
	"log"
	"os"
	"s3policy/internal/config"
	"s3policy/internal/service"

	"github.com/alexflint/go-arg"
)

var version string

type commonParams struct {
	Bucket string `arg:"-b,required" help:"Bucket name"`
}

type PutPolicyCmd struct {
	PolicyFile string `arg:"-f,--file,required" help:"File with policy"`
	commonParams
}

type GetPolicyCmd struct {
	commonParams
}

type DeletePolicyCmd struct {
	commonParams
}

type ListBucketsCmd struct {
}

type CreateBucketCmd struct {
	commonParams
}

type DeleteBucketCmd struct {
	commonParams
}

type GetBucketACLCmd struct {
	commonParams
	Verbose bool `arg:"-v,--verbose" help:"Display extended info"`
}

type GetBucketCORSCmd struct {
	commonParams
}

type PutCORSCmd struct {
	commonParams
}

type DeleteCORSCmd struct {
	commonParams
}

type args struct {
	ConfigFile string `arg:"-c, --config" help:"Config INI file" default:"./config.ini"`
	Profile    string `arg:"-p,required" help:"Profile name"`

	PutPolicy     *PutPolicyCmd     `arg:"subcommand:put-policy" help:"Put a policy"`
	GetPolicy     *GetPolicyCmd     `arg:"subcommand:get-policy" help:"Display a policy"`
	DeletePolicy  *DeletePolicyCmd  `arg:"subcommand:del-policy" help:"Delete a policy"`
	ListBuckets   *ListBucketsCmd   `arg:"subcommand:list-buckets" help:"List buckets"`
	CreateBucket  *CreateBucketCmd  `arg:"subcommand:create-bucket" help:"Create a bucket"`
	DeleteBucket  *DeleteBucketCmd  `arg:"subcommand:del-bucket" help:"Delete a bucket"`
	GetBucketACL  *GetBucketACLCmd  `arg:"subcommand:get-bucket-acl" help:"Get bucket's ACL"`
	GetBucketCORS *GetBucketCORSCmd `arg:"subcommand:get-cors" help:"Get bucket's CORS"`
	PutCORS       *PutCORSCmd       `arg:"subcommand:put-cors" help:"Put bucket's CORS"`
	DeleteCORS    *DeleteCORSCmd    `arg:"subcommand:del-cors" help:"Delete bucket's CORS"`
}

func (args) Version() string {
	return fmt.Sprintf("Version %s", version)
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("ERROR happend: %s", err)
			// debug.PrintStack()
			os.Exit(1)
		}
	}()

	var args args
	parser := arg.MustParse(&args)

	cfg := config.NewConfig(args.ConfigFile)
	client := service.NewService(cfg, args.Profile)

	var err error

	switch {

	case args.GetPolicy != nil:
		err = client.GetPolicy(args.GetPolicy.Bucket)

	case args.PutPolicy != nil:
		err = client.PutPolicy(args.PutPolicy.Bucket, args.PutPolicy.PolicyFile)

	case args.DeletePolicy != nil:
		err = client.DeletePolicy(args.DeletePolicy.Bucket)

	case args.ListBuckets != nil:
		err = client.ListBuckets()

	case args.CreateBucket != nil:
		err = client.CreateBucket(args.CreateBucket.Bucket)

	case args.DeleteBucket != nil:
		err = client.DeleteBucket(args.DeleteBucket.Bucket)

	case args.GetBucketACL != nil:
		err = client.GetBucketACL(args.GetBucketACL.Bucket, args.GetBucketACL.Verbose)

	case args.GetBucketCORS != nil:
		err = client.GetBucketCORS(args.GetBucketCORS.Bucket)

	case args.PutCORS != nil:
		err = client.PutBucketCORS(args.PutCORS.Bucket)

	case args.DeleteCORS != nil:
		err = client.DeleteCORS(args.DeleteCORS.Bucket)

	default:

		parser.WriteHelp(os.Stdin)
		panic("Undefined command")
	}

	if err != nil {
		panic(err)
	}
}
