Tool for managing S3
---------------------------

Lite tool for managing S3 buckets and policies. It does not implement all the actions that can be performed with buckets, only creating\deleting\viewing buckets and creating\deleting\viewing policies.

Configurations
---------------------------
Default configuration file is __./config.ini__. It may have many profiles:

```
[<profile name>]
id = "<Access Key ID>"
key = "<Secret Key>"
region = "<region name>"
endpoint = "https://<url to service>"

[<profile name2>]
...
...

```

How to build
---------------------------

For Linux
```
GOOS=linux go build -ldflags "-X main.version=X.X.X" -o s3policy ./cmd/policy/
```

For Windows
```
GOOS=windows go build -ldflags "-X main.version=X.X.X" -o s3policy ./cmd/policy/
```

Example
---------------------------

To get a policy for bucket __test__:

```
s3policy get-policy -p dev -b ctest
```

It returns something like that
```
{
        "Id": "Policy11111111",
        "Version": "2012-10-17",
        "Statement": [
          {
                "Sid": "Stmt222222",
                "Action": [
                  "s3:GetObject"
                ],
                "Effect": "Allow",
                "Resource": "arn:aws:s3:::test/*",
                "Principal": "*"
          }
        ]
  }
```

To get a policy for bucket __test__:

```
s3policy put-policy -p dev -b ctest -f ./test-plolicy.json
```

You should provide a policy in the file __test-plolicy.json__
