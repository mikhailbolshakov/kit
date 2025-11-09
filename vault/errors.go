package vault

import (
	"context"
	"github.com/mikhailbolshakov/kit"
)

const (
	ErrCodeVaultNewClient               = "VLT-001"
	ErrCodeVaultClientWrite             = "VLT-002"
	ErrCodeVaultClientEmptyResponseData = "VLT-003"
	ErrCodeVaultUnmarshalResponse       = "VLT-004"
	ErrCodeVaultUnmountPKIEngine        = "VLT-005"
	ErrCodeVaultClientRead              = "VLT-006"

	ErrCodeCertificateParse                   = "CRT-001"
	ErrCodeCertificateParseCSR                = "CRT-002"
	ErrCodeCertificateCSRNotValid             = "CRT-003"
	ErrCodeCertificateUnrecognizedBlockType   = "CRT-004"
	ErrCodeCertificateUnsupportedKeyAlgorithm = "CRT-005"
	ErrCodeCertificateUnsupportedKeyType      = "CRT-006"
	ErrCodeCertificatePemEncode               = "CRT-007"
	ErrCodeCertificateNotValid                = "CRT-008"
	ErrCodeCertificateBlockEmpty              = "CRT-009"
	ErrCodeCertificateInvalidBlockType        = "CRT-010"
)

var (
	ErrVaultNewClient = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeVaultNewClient, "vault: client init failed").Wrap(cause).C(ctx).Err()
	}
	ErrVaultClientWrite = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeVaultClientWrite, "vault: client logical write failed").Wrap(cause).C(ctx).Err()
	}
	ErrVaultClientRead = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeVaultClientRead, "vault: client logical read failed").Wrap(cause).C(ctx).Err()
	}
	ErrCertificatePemEncode = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeCertificatePemEncode, "certificate: pem encode failed").Wrap(cause).C(ctx).Err()
	}
	ErrVaultClientEmptyResponseData = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeVaultClientEmptyResponseData, "vault: empty response").Business().C(ctx).Err()
	}
	ErrCertificateUnsupportedKeyType = func(ctx context.Context, algorithm int) error {
		return kit.NewAppErrBuilder(ErrCodeCertificateUnsupportedKeyType, "certificate: unsupported key type").F(kit.KV{"algorithm": algorithm}).Business().C(ctx).Err()
	}
	ErrCertificateUnsupportedKeyAlgorithm = func(ctx context.Context, algorithm int) error {
		return kit.NewAppErrBuilder(ErrCodeCertificateUnsupportedKeyAlgorithm, "certificate: unsupported key algorithm").F(kit.KV{"algorithm": algorithm}).Business().C(ctx).Err()
	}
	ErrVaultUnmarshalResponse = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeVaultUnmarshalResponse, "vault: unmarshal response failed").Wrap(cause).C(ctx).Err()
	}
	ErrVaultUnmountPKIEngine = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeVaultUnmountPKIEngine, "vault: unmount pki engine failed").Wrap(cause).C(ctx).Err()
	}
	ErrCertificateNotValid = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeCertificateNotValid, "certificate: not valid").Business().Wrap(cause).C(ctx).Err()
	}
	ErrCertificateParseCertificate = func(ctx context.Context, cause error) error {
		return kit.NewAppErrBuilder(ErrCodeCertificateParse, "certificate: parse certificate failed").Wrap(cause).C(ctx).Err()
	}
	ErrCertificateBlockEmpty = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCertificateBlockEmpty, "certificate: empty block").Business().C(ctx).Err()
	}
	ErrCertificateParseCSR = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCertificateParseCSR, "certificate: parse CSR failed").Business().C(ctx).Err()
	}
	ErrCertificateCSRNotValid = func(ctx context.Context) error {
		return kit.NewAppErrBuilder(ErrCodeCertificateCSRNotValid, "certificate: CSR not valid").Business().C(ctx).Err()
	}
	ErrCertificateUnrecognizedBlockType = func(ctx context.Context, tp string) error {
		return kit.NewAppErrBuilder(ErrCodeCertificateUnrecognizedBlockType, "certificate: certificate block type is unrecognized").F(kit.KV{"type": tp}).Business().C(ctx).Err()
	}
	ErrCertificateInvalidBlockType = func(ctx context.Context, tp string) error {
		return kit.NewAppErrBuilder(ErrCodeCertificateInvalidBlockType, "certificate: certificate block type is invalid").F(kit.KV{"type": tp}).Business().C(ctx).Err()
	}
)
