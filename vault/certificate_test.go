package vault

import (
	"github.com/mikhailbolshakov/kit"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type certificateSuite struct {
	kit.Suite
}

func (s *certificateSuite) SetupSuite() {
	s.Suite.Init(func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) })
}

func TestCertificateSuite(t *testing.T) {
	suite.Run(t, new(certificateSuite))
}

func (s *certificateSuite) SetupTest() {
}

func (s *certificateSuite) Test_ParseVerify() {
	resp, err := http.Get("https://www.google.com")
	s.NoError(err)
	certs := resp.TLS.PeerCertificates
	s.NotEmpty(certs)

	s.NoError(VerifyCertificate(s.Ctx, certs[0], certs[1:]...))
	s.NoError(VerifyCertificate(s.Ctx, certs[2]))
	s.NoError(VerifyCertificate(s.Ctx, certs[1]))
	// without ca can't be verified
	s.AssertAppErr(VerifyCertificate(s.Ctx, certs[0]), ErrCodeCertificateNotValid)

	chainStr, err := ToChainPemString(s.Ctx, certs)
	s.NoError(err)
	s.NotEmpty(chainStr)

	certStr, err := ToCertificatePemString(s.Ctx, certs[0])
	s.NoError(err)
	s.NotEmpty(certStr)
	certStr, err = ToCertificatePemString(s.Ctx, certs[1])
	s.NoError(err)
	s.NotEmpty(certStr)
	certStr, err = ToCertificatePemString(s.Ctx, certs[2])
	s.NoError(err)
	s.NotEmpty(certStr)

	chain, err := ParseCertificateChain(s.Ctx, chainStr)
	s.NoError(err)
	s.NotEmpty(chain)

	s.NoError(VerifyCertificate(s.Ctx, chain[0], chain[1:]...))
	s.NoError(VerifyCertificate(s.Ctx, chain[2]))
	s.NoError(VerifyCertificate(s.Ctx, chain[1]))
	// without ca can't be verified
	s.AssertAppErr(VerifyCertificate(s.Ctx, chain[0]), ErrCodeCertificateNotValid)

	parsed, err := ToCertificatePemString(s.Ctx, chain[0])
	s.NoError(err)
	s.NotEmpty(parsed)
	parsed, err = ToCertificatePemString(s.Ctx, chain[1])
	s.NoError(err)
	s.NotEmpty(parsed)
	parsed, err = ToCertificatePemString(s.Ctx, chain[2])
	s.NoError(err)
	s.NotEmpty(parsed)

}
