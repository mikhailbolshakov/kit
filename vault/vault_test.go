//go:build integration

package vault

import (
	"crypto/tls"
	"github.com/mikhailbolshakov/kit"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type vaultSuite struct {
	kit.Suite
	logger kit.CLoggerFunc
	vault  Vault
}

func (s *vaultSuite) SetupSuite() {
	s.logger = func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) }
	s.Suite.Init(s.logger)

	s.vault = NewVault(s.logger).Path("pki")
	s.NoError(s.vault.Init(s.Ctx, &Config{
		Mock:       false,
		Url:        "http://localhost:8201",
		Token:      "00000000-0000-0000-0000-000000000000",
		TimeoutSec: 10,
	}))

	_ = s.vault.UnmountPKIEngine(s.Ctx)
	err := s.vault.EnablePKIEngine(s.Ctx)
	s.NoError(err)
}

func (s *vaultSuite) TearDownSuite() {
	s.vault.UnmountPKIEngine(s.Ctx)
}

func TestVaultSuite(t *testing.T) {
	suite.Run(t, new(vaultSuite))
}

func (s *vaultSuite) SetupTest() {}

func (s *vaultSuite) Test_CreateByIntermediate() {
	path := kit.NewRandString()
	vaultRoot := NewVault(s.logger).Path("root_pki_" + path)
	s.NoError(vaultRoot.Init(s.Ctx, &Config{
		Url:        "http://localhost:8201",
		Token:      "00000000-0000-0000-0000-000000000000",
		TimeoutSec: 10,
	}))
	s.NoError(vaultRoot.EnablePKIEngine(s.Ctx))
	defer vaultRoot.UnmountPKIEngine(s.Ctx)

	vaultInter := NewVault(s.logger).Path("inter_pki_" + path)
	s.NoError(vaultInter.Init(s.Ctx, &Config{
		Url:        "http://localhost:8201",
		Token:      "00000000-0000-0000-0000-000000000000",
		TimeoutSec: 10,
	}))
	s.NoError(vaultInter.EnablePKIEngine(s.Ctx))
	defer vaultInter.UnmountPKIEngine(s.Ctx)

	vaultLeaf := NewVault(s.logger).Path("leaf_pki_" + path)
	s.NoError(vaultLeaf.Init(s.Ctx, &Config{
		Url:        "http://localhost:8201",
		Token:      "00000000-0000-0000-0000-000000000000",
		TimeoutSec: 10,
	}))
	s.NoError(vaultLeaf.EnablePKIEngine(s.Ctx))
	defer vaultLeaf.UnmountPKIEngine(s.Ctx)

	rootName := kit.NewRandString() + ".com"
	root, err := vaultRoot.GenerateRootCA(s.Ctx, &CertificateRequest{CommonName: rootName})
	s.NoError(err)
	s.NotEmpty(root)

	csr, err := vaultInter.GenerateCSR(s.Ctx, &CSRRequest{
		Type: CSRTypeInternal,
		CertificateRequest: CertificateRequest{
			CommonName:         rootName,
			Country:            []string{"EG", "AU"},
			Province:           []string{"State1", "State2"},
			Locality:           []string{"City"},
			Organization:       []string{"Company Ltd"},
			OrganizationalUnit: []string{"IT"},
		},
	})
	s.NoError(err)
	s.NotEmpty(csr)

	cert, err := vaultRoot.SignIntermediateCA(s.Ctx, &SignRequest{
		CSR:      csr.CSR,
		Format:   SignFormatPemBundle,
		NotAfter: ToDateFormat(kit.Now().Add(time.Hour)),
	})
	s.NoError(err)
	s.NotEmpty(cert)

	signed, err := vaultInter.SetSignedCertificate(s.Ctx, &SetSignedCertificateRequest{Certificate: cert.Certificate})
	s.NoError(err)
	s.NotEmpty(signed)

	roleName := kit.NewRandString()
	role, err := vaultInter.CreateRole(s.Ctx, &RoleRequest{Name: roleName, AllowSubdomains: kit.BoolPtr(true), AllowedDomains: []string{rootName}})
	s.NoError(err)
	s.NotEmpty(role)

	csr, err = vaultLeaf.GenerateCSR(s.Ctx, &CSRRequest{
		Type: CSRTypeExported,
		CertificateRequest: CertificateRequest{
			CommonName:         "csr1." + rootName,
			Country:            []string{"EG", "AU"},
			Province:           []string{"State1", "State2"},
			Locality:           []string{"City"},
			Organization:       []string{"Company Ltd"},
			OrganizationalUnit: []string{"IT"},
		},
	})
	s.NoError(err)
	s.NotEmpty(csr)

	cert, err = vaultInter.Sign(s.Ctx, roleName, &SignRequest{
		CSR:      csr.CSR,
		Format:   SignFormatPemBundle,
		NotAfter: ToDateFormat(kit.Now().Add(time.Hour)),
	})
	s.NoError(err)
	s.NotEmpty(cert)
}

func (s *vaultSuite) Test_SignByRole() {
	rootName := kit.NewRandString() + ".com"
	root, err := s.vault.GenerateRootCA(s.Ctx, &CertificateRequest{CommonName: rootName})
	s.NoError(err)
	s.NotEmpty(root)

	roleName := kit.NewRandString()
	role, err := s.vault.CreateRole(s.Ctx, &RoleRequest{Name: roleName, AllowSubdomains: kit.BoolPtr(true), AllowedDomains: []string{rootName}})
	s.NoError(err)
	s.NotEmpty(role)

	csr, err := s.vault.GenerateCSR(s.Ctx, &CSRRequest{
		Type: CSRTypeInternal,
		CertificateRequest: CertificateRequest{
			CommonName:         "csr1." + rootName,
			Country:            []string{"EG", "AU"},
			Province:           []string{"State1", "State2"},
			Locality:           []string{"City"},
			Organization:       []string{"Company Ltd"},
			OrganizationalUnit: []string{"IT"},
		},
	})
	s.NoError(err)
	s.NotEmpty(csr)

	cert, err := s.vault.Sign(s.Ctx, roleName, &SignRequest{
		CSR:      csr.CSR,
		Format:   SignFormatPemBundle,
		NotAfter: ToDateFormat(kit.Now().Add(time.Hour)),
	})
	s.NoError(err)
	s.NotEmpty(cert)

}

func (s *vaultSuite) Test_CreateVerifyIssued() {
	root, err := s.vault.GenerateRootCA(s.Ctx, &CertificateRequest{CommonName: kit.NewRandString()})
	s.NoError(err)
	s.NotEmpty(root)

	csr, err := s.vault.GenerateCSR(s.Ctx, &CSRRequest{
		Type: CSRTypeInternal,
		CertificateRequest: CertificateRequest{
			CommonName:         "example.com",
			Country:            []string{"EG", "AU"},
			Province:           []string{"State1", "State2"},
			Locality:           []string{"City"},
			Organization:       []string{"Company Ltd"},
			OrganizationalUnit: []string{"IT"},
		},
	})
	s.NoError(err)
	s.NotEmpty(csr)

	cert, err := s.vault.SignIntermediateCA(s.Ctx, &SignRequest{
		CSR:      csr.CSR,
		Format:   SignFormatPemBundle,
		NotAfter: ToDateFormat(kit.Now().Add(time.Hour)),
	})
	s.NoError(err)
	s.NotEmpty(cert)

	parsedCert, err := ParseCertificateChain(s.Ctx, cert.Certificate)
	s.NoError(err)
	parsedCSR, err := ParseCSR(s.Ctx, csr.CSR)
	s.NoError(err)

	s.NoError(VerifyCertificate(s.Ctx, parsedCert[0], parsedCert[1:]...))

	// check if issued by csr
	issued, err := IssuedByCSR(s.Ctx, parsedCSR, parsedCert[0])
	s.NoError(err)
	s.True(issued)
}

func (s *vaultSuite) Test_CreateCSR() {
	root, err := s.vault.GenerateRootCA(s.Ctx, &CertificateRequest{CommonName: kit.NewRandString()})
	s.NoError(err)
	s.NotEmpty(root)

	name := kit.NewRandString()
	role, err := s.vault.CreateRole(s.Ctx, &RoleRequest{Name: name})
	s.NoError(err)
	s.NotEmpty(role)

	csr, err := s.vault.GenerateCSR(s.Ctx, &CSRRequest{
		Type: CSRTypeInternal,
		CertificateRequest: CertificateRequest{
			CommonName:         "example.com",
			Country:            []string{"EG", "AU"},
			Province:           []string{"State1", "State2"},
			Locality:           []string{"City"},
			Organization:       []string{"Company Ltd"},
			OrganizationalUnit: []string{"IT"},
		},
	})
	s.NoError(err)
	s.NotEmpty(csr)

	cert, err := s.vault.SignIntermediateCA(s.Ctx, &SignRequest{
		CSR:      csr.CSR,
		Format:   SignFormatPemBundle,
		NotAfter: ToDateFormat(kit.Now().Add(time.Hour)),
		Name:     name,
	})
	s.NoError(err)
	s.NotEmpty(cert)

	cert, err = s.vault.SignIntermediateCA(s.Ctx, &SignRequest{
		CSR:    csr.CSR,
		Format: SignFormatPemBundle,
		Name:   name,
		TTL:    ToTTLHours(24),
	})
	s.NoError(err)
	s.NotEmpty(cert)

	signed, err := s.vault.SetSignedCertificate(s.Ctx, &SetSignedCertificateRequest{Certificate: cert.Certificate})
	s.NoError(err)
	s.NotEmpty(signed)

	got, err := s.vault.GetCertificate(s.Ctx, cert.SerialNumber)
	s.NoError(err)
	s.NotEmpty(cert)
	s.Contains(cert.Certificate, got.Certificate)

	parsed, err := ParseCertificateChain(s.Ctx, cert.Certificate)
	s.NoError(err)
	s.Len(parsed, 2)

	s.NoError(VerifyCertificate(s.Ctx, parsed[0], parsed[1:]...))

	// to objects
	csrX509, err := ParseCSR(s.Ctx, csr.CSR)
	s.NoError(err)
	s.NotEmpty(csrX509)

	rootX509, err := ParseCertificate(s.Ctx, root.Certificate)
	s.NoError(err)
	s.NotEmpty(rootX509)

	certX509, err := ParseCertificate(s.Ctx, cert.Certificate)
	s.NoError(err)
	s.NotEmpty(certX509)

	gotX509, err := ParseCertificate(s.Ctx, got.Certificate)
	s.NoError(err)
	s.NotEmpty(gotX509)

	// check if issued by csr
	issued, err := IssuedByCSR(s.Ctx, csrX509, gotX509)
	s.NoError(err)
	s.True(issued)

	issued, err = IssuedByCSR(s.Ctx, csrX509, rootX509)
	s.NoError(err)
	s.False(issued)

	issued, err = IssuedByCSR(s.Ctx, csrX509, certX509)
	s.NoError(err)
	s.True(issued)

	csrStr, err := ToCSRPemString(s.Ctx, csrX509)
	s.NoError(err)
	s.NotEmpty(csrStr)

	rootStr, err := ToCertificatePemString(s.Ctx, rootX509)
	s.NoError(err)
	s.NotEmpty(rootStr)

	certStr, err := ToCertificatePemString(s.Ctx, certX509)
	s.NoError(err)
	s.NotEmpty(certStr)
}

func (s *vaultSuite) Test_IssueByRole_TlsCertificate() {

	path := kit.NewRandString()
	rootName := kit.NewRandString() + ".com"

	vaultRoot := NewVault(s.logger).Path("root_pki_" + path)
	s.NoError(vaultRoot.Init(s.Ctx, &Config{
		Url:        "http://localhost:8201",
		Token:      "00000000-0000-0000-0000-000000000000",
		TimeoutSec: 10,
	}))
	s.NoError(vaultRoot.EnablePKIEngine(s.Ctx))
	defer vaultRoot.UnmountPKIEngine(s.Ctx)

	root, err := vaultRoot.GenerateRootCA(s.Ctx, &CertificateRequest{CommonName: kit.NewRandString()})
	s.NoError(err)
	s.NotEmpty(root)

	vaultInter := NewVault(s.logger).Path("inter_pki_" + path)
	s.NoError(vaultInter.Init(s.Ctx, &Config{
		Url:        "http://localhost:8201",
		Token:      "00000000-0000-0000-0000-000000000000",
		TimeoutSec: 10,
	}))
	s.NoError(vaultInter.EnablePKIEngine(s.Ctx))
	defer vaultInter.UnmountPKIEngine(s.Ctx)

	csr, err := vaultInter.GenerateCSR(s.Ctx, &CSRRequest{
		Type: CSRTypeInternal,
		CertificateRequest: CertificateRequest{
			CommonName:         rootName,
			Country:            []string{"EG", "AU"},
			Province:           []string{"State1", "State2"},
			Locality:           []string{"City"},
			Organization:       []string{"Company Ltd"},
			OrganizationalUnit: []string{"IT"},
		},
	})
	s.NoError(err)
	s.NotEmpty(csr)

	cert, err := vaultRoot.SignIntermediateCA(s.Ctx, &SignRequest{
		CSR:      csr.CSR,
		Format:   SignFormatPemBundle,
		NotAfter: ToDateFormat(kit.Now().Add(time.Hour * 2)),
	})
	s.NoError(err)
	s.NotEmpty(cert)

	signed, err := vaultInter.SetSignedCertificate(s.Ctx, &SetSignedCertificateRequest{Certificate: cert.Certificate})
	s.NoError(err)
	s.NotEmpty(signed)

	signed, err = vaultInter.SetSignedCertificate(s.Ctx, &SetSignedCertificateRequest{Certificate: root.Certificate})
	s.NoError(err)
	s.NotEmpty(signed)

	intChain, err := vaultInter.GetChain(s.Ctx)
	s.NoError(err)
	s.NotEmpty(intChain)

	rootChain, err := vaultRoot.GetChain(s.Ctx)
	s.NoError(err)
	s.NotEmpty(rootChain)

	roleName := kit.NewRandString()
	role, err := vaultInter.CreateRole(s.Ctx, &RoleRequest{Name: roleName, AllowSubdomains: kit.BoolPtr(true), AllowedDomains: []string{rootName}})
	s.NoError(err)
	s.NotEmpty(role)

	// issue leaf certificate by intermediate
	leaf, err := vaultInter.Issue(s.Ctx, roleName, &IssueRequest{
		CommonName: "test." + rootName,
		Format:     SignFormatPemBundle,
		NotAfter:   ToDateFormat(kit.Now().Add(time.Hour)),
	})
	s.NoError(err)
	s.NotEmpty(leaf)

	tlsCert, err := tls.X509KeyPair([]byte(leaf.Certificate), []byte(leaf.PrivateKey))
	s.NoError(err)
	s.NotEmpty(tlsCert)

	got, err := vaultInter.GetCertificate(s.Ctx, leaf.SerialNumber)
	s.NoError(err)
	s.NotEmpty(got)
	s.Contains(leaf.Certificate, got.Certificate)

}
