package response

type Ec2 struct {
	Ec2Name      string `json:"ec2_name"`
	InstanceSize string `json:"instance_size"`
}

type Vpc struct {
	CidrRange      string `json:"cidr_range"`
	InstanceNumber string `json:"instance_number"`
	VpcID          string `json:"vpc_id"`
}

type ElasticCache struct {
	InstanceSize string `json:"instance_size"`
}

type Rds struct {
	InstanceSize string `json:"instance_size"`
}

type Lambda []struct {
	FunctionName string `json:"function_name"`
}

type CloudFront struct {
	Link string `json:"link"`
}

type S3Bucket []struct {
	Name string `json:"name"`
}

type Route53 struct {
	Link       string `json:"link"`
	HostedZone string `json:"hosted_zone"`
}
type AwsS3Response struct {
	Ec2            Ec2          `json:"ec2"`
	Vpc            Vpc          `json:"vpc"`
	SslCertificate string       `json:"ssl_certificate"`
	ElasticCache   ElasticCache `json:"elastic_cache"`
	Rds            Rds          `json:"rds"`
	Lambda         interface{}  `json:"lambda"`
	CloudFront     CloudFront   `json:"cloud_front"`
	S3Bucket       S3Bucket     `json:"s3_bucket"`
	Route53        Route53      `json:"route53"`
}
