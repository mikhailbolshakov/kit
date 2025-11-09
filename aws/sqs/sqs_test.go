//go:build dev

package sqs

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/suite"
	"gitlab.com/algmib/kit"
	kitAws "gitlab.com/algmib/kit/aws"
)

type s3TestSuite struct {
	kit.Suite
	logger kit.CLoggerFunc
}

func (s *s3TestSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)
}

func TestS3Suite(t *testing.T) {
	suite.Run(t, new(s3TestSuite))
}

var (
	awsCfg = &kitAws.Config{
		Region:              "eu-central-1",
		AccessKeyId:         "access_key_id",
		SecretAccessKey:     "secret_access_key",
		SharedConfigProfile: "algmib/dev",
	}
)

func (s *s3TestSuite) Test_Init() {
	// init client
	client := NewClient(awsCfg, s.logger)
	s.NoError(client.Init(s.Ctx))
	s.NotEmpty(client.sqsClient)

	_, err := client.GetQueueURL(s.Ctx, &sqs.GetQueueUrlInput{
		QueueName:              kit.StringPtr("ext-storage-dev"),
		QueueOwnerAWSAccountId: nil,
	})
	s.NoError(err)
}
