//go:build dev

package ftp

import (
	_ "github.com/jlaffaye/ftp"
	"github.com/stretchr/testify/suite"
	"gitlab.com/algmib/kit"
	"testing"
)

type ftpTestSuite struct {
	kit.Suite
}

func (s *ftpTestSuite) SetupSuite() {
	s.Suite.Init(func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) })
}

func (s *ftpTestSuite) SetupTest() {}

func TestFtpSuite(t *testing.T) {
	suite.Run(t, new(ftpTestSuite))
}

func (s *ftpTestSuite) Test_Read() {
	//buf, err := ReadFileByUrl(s.Ctx, "ftp://firmware:firmware@localhost:8021/firmware/123/456/article.pdf")
	buf, err := ReadFileByUrl(s.Ctx, "ftp://firmware:4BSaqRsIJeYbaQR@minio.algmib.tech:8021/firmware/d1b658a1-495a-4a70-8458-45e2d154e4a1/927133ce-f636-4952-a02a-b5fe41ab874a/20240306125651_firmware_j6e13g.bin")
	s.NoError(err)
	s.NotEmpty(buf)
}

func (s *ftpTestSuite) Test_Load() {
	//buf, err := ReadFileByUrl(s.Ctx, "ftp://firmware:firmware@localhost:8021/firmware/123/456/article.pdf")
	err := LoadFileByUrl(s.Ctx, "ftp://firmware:4BSaqRsIJeYbaQR@minio.algmib.tech:8021/firmware/d1b658a1-495a-4a70-8458-45e2d154e4a1/927133ce-f636-4952-a02a-b5fe41ab874a/20240306125651_firmware_j6e134.bin", []byte{111, 222, 111, 222})
	s.NoError(err)
}
