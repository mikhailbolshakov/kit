package cluster

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/stretchr/testify/suite"
	"testing"
)

type serviceUtilsTestSuite struct {
	kit.Suite
}

func (s *serviceUtilsTestSuite) SetupSuite() {
	s.Suite.Init(nil)
}

func TestServiceUtilsSuite(t *testing.T) {
	suite.Run(t, new(serviceUtilsTestSuite))
}

func (s *serviceUtilsTestSuite) Test_GetServiceRootPath_WhenEmptyInput() {
	s.Empty(GetServiceRootPath(""))
}

func (s *serviceUtilsTestSuite) Test_GetServiceRootPath_WhenNotExistent() {
	s.Empty(GetServiceRootPath("some"))
}

func (s *serviceUtilsTestSuite) Test_GetServiceRootPath_WhenKit() {
	s.NotEmpty(GetServiceRootPath("kit"))
}
