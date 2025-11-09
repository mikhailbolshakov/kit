package vault

import (
	"bytes"
	"context"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"gitlab.com/algmib/kit"
)

func ToCSRPemString(ctx context.Context, csr *x509.CertificateRequest) (string, error) {
	return toString(ctx, []*pem.Block{{
		Type:  CsrPemBlockType,
		Bytes: csr.Raw,
	}})
}

func ToCSRRawPemString(ctx context.Context, raw []byte) (string, error) {
	return toString(ctx, []*pem.Block{{
		Type:  CsrPemBlockType,
		Bytes: raw,
	}})
}

func ToCertificatePemString(ctx context.Context, cert *x509.Certificate) (string, error) {
	return toString(ctx, []*pem.Block{{
		Type:  CertificatePemBlockType,
		Bytes: cert.Raw,
	}})
}

func ToCertificateRawPemString(ctx context.Context, raw []byte) (string, error) {
	return toString(ctx, []*pem.Block{{
		Type:  CertificatePemBlockType,
		Bytes: raw,
	}})
}

func ToPrivateKeyRsaPemString(ctx context.Context, privateKey *rsa.PrivateKey) (string, error) {
	return toString(ctx, []*pem.Block{{
		Type:  fmt.Sprintf("%s %s", RSA, PrivateKeyPemBlockType),
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}})
}

func ToChainPemString(ctx context.Context, certs []*x509.Certificate) (string, error) {
	return toString(ctx, kit.ConvertSlice(certs, func(cert *x509.Certificate) *pem.Block {
		return &pem.Block{
			Type:  CertificatePemBlockType,
			Bytes: cert.Raw,
		}
	}))
}

func ParseCertificateChain(ctx context.Context, chain string) ([]*x509.Certificate, error) {
	var res []*x509.Certificate
	for {
		block, rest := pem.Decode([]byte(chain))
		if block == nil {
			break // no more PEM data
		}
		if block.Type != CertificatePemBlockType {
			return nil, ErrCertificateInvalidBlockType(ctx, block.Type)
		}
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, ErrCertificateParseCertificate(ctx, err)
		}
		// append data
		res = append(res, cert)
		// update for next loop iteration
		chain = string(rest)
	}
	return res, nil
}

func ParseCertificate(ctx context.Context, certificate string) (*x509.Certificate, error) {
	block, rest := pem.Decode([]byte(certificate))
	if block == nil {
		return nil, ErrCertificateBlockEmpty(ctx)
	}
	if block.Type == CertificatePemBlockType {
		c, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, ErrCertificateParseCertificate(ctx, err)
		}
		return c, nil
	}
	if strings.Contains(block.Type, PrivateKeyPemBlockType) {
		// if the first block is a private key, it means we got a bundle.
		// parse the rest of the chain
		return ParseCertificate(ctx, string(rest))
	}
	return nil, ErrCertificateUnrecognizedBlockType(ctx, block.Type)
}

func ParseCSR(ctx context.Context, csr string) (*x509.CertificateRequest, error) {
	block, _ := pem.Decode([]byte(csr))
	if block == nil {
		return nil, ErrCertificateParseCSR(ctx)
	}

	if block.Type != CsrPemBlockType {
		return nil, ErrCertificateCSRNotValid(ctx)
	}

	rq, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rq, nil
}

func VerifyCertificate(ctx context.Context, cert *x509.Certificate, caCerts ...*x509.Certificate) error {
	opt := x509.VerifyOptions{}

	if len(caCerts) != 0 {
		roots := x509.NewCertPool()
		// load root and intermediate CAs
		for _, caCert := range caCerts {
			roots.AddCert(caCert)
		}
		opt.Roots = roots
	}
	// Verify the certificate
	if _, err := cert.Verify(opt); err != nil {
		return ErrCertificateNotValid(ctx, err)
	}
	return nil
}

func IssuedByCSR(ctx context.Context, csr *x509.CertificateRequest, certificate *x509.Certificate) (bool, error) {
	// compare algorithm
	if csr.PublicKeyAlgorithm == x509.UnknownPublicKeyAlgorithm {
		return false, ErrCertificateUnsupportedKeyAlgorithm(ctx, int(csr.PublicKeyAlgorithm))
	}
	if csr.PublicKeyAlgorithm != certificate.PublicKeyAlgorithm {
		return false, nil
	}
	// compare public keys
	return comparePublicKeys(ctx, csr.PublicKeyAlgorithm, csr.PublicKey, certificate.PublicKey)
}

func toString(ctx context.Context, blocks []*pem.Block) (string, error) {
	var b = &bytes.Buffer{}
	for _, block := range blocks {
		if err := pem.Encode(b, block); err != nil {
			return "", ErrCertificatePemEncode(ctx, err)
		}
	}
	return b.String(), nil
}

func comparePublicKeys(ctx context.Context, algorithm x509.PublicKeyAlgorithm, key1, key2 any) (bool, error) {
	switch algorithm {
	case x509.RSA:
		rsa1, ok := key1.(*rsa.PublicKey)
		if !ok {
			return false, ErrCertificateUnsupportedKeyType(ctx, int(algorithm))
		}
		rsa2, ok := key2.(*rsa.PublicKey)
		if !ok {
			return false, ErrCertificateUnsupportedKeyType(ctx, int(algorithm))
		}
		return rsa1.N.Cmp(rsa2.N) == 0 && rsa1.E == rsa2.E, nil
	case x509.DSA:
		dsa1, ok := key1.(*dsa.PublicKey)
		if !ok {
			return false, ErrCertificateUnsupportedKeyType(ctx, int(algorithm))
		}
		dsa2, ok := key2.(*dsa.PublicKey)
		if !ok {
			return false, ErrCertificateUnsupportedKeyType(ctx, int(algorithm))
		}
		return dsa1.Y.Cmp(dsa2.Y) == 0 && dsa1.P.Cmp(dsa2.P) == 0 &&
			dsa1.Q.Cmp(dsa2.Q) == 0 && dsa1.G.Cmp(dsa2.G) == 0, nil
	case x509.ECDSA:
		ecdsa1, ok := key1.(*ecdsa.PublicKey)
		if !ok {
			return false, ErrCertificateUnsupportedKeyType(ctx, int(algorithm))
		}
		ecdsa2, ok := key2.(*ecdsa.PublicKey)
		if !ok {
			return false, ErrCertificateUnsupportedKeyType(ctx, int(algorithm))
		}
		return ecdsa1.X.Cmp(ecdsa2.X) == 0 && ecdsa1.Y.Cmp(ecdsa2.Y) == 0 && ecdsa1.Curve == ecdsa2.Curve, nil
	case x509.Ed25519:
		ed255191, ok := key1.(*ed25519.PublicKey)
		if !ok {
			return false, ErrCertificateUnsupportedKeyType(ctx, int(algorithm))
		}
		ed255192, ok := key2.(*ed25519.PublicKey)
		if !ok {
			return false, ErrCertificateUnsupportedKeyType(ctx, int(algorithm))
		}
		return bytes.Equal(*ed255191, *ed255192), nil
	default:
		return false, ErrCertificateUnsupportedKeyAlgorithm(ctx, int(algorithm))

	}
}
