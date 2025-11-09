package vault

import (
	"context"
	"fmt"
	"time"

	lib "github.com/hashicorp/vault/api"
	"github.com/mikhailbolshakov/kit"
)

const (
	PathPki   = "pki"
	EnginePki = "pki"

	MaxTTL = "720h"
	TTL    = "24h"

	CsrPemBlockType         = "CERTIFICATE REQUEST"
	CertificatePemBlockType = "CERTIFICATE"
	PrivateKeyPemBlockType  = "PRIVATE KEY"

	RSA     KeyAlgorithm = "RSA"
	DSA     KeyAlgorithm = "DSA" // Only supported for parsing.
	ECDSA   KeyAlgorithm = "ECDSA"
	Ed25519 KeyAlgorithm = "Ed25519"
)

type KeyAlgorithm string

type Vault interface {
	Path(path string) Vault
	Init(ctx context.Context, cfg *Config) error

	GenerateCSR(ctx context.Context, rq *CSRRequest) (*CSR, error)
	GenerateRootCA(ctx context.Context, rq *CertificateRequest) (*Certificate, error)
	SignIntermediateCA(ctx context.Context, rq *SignRequest) (*Certificate, error)
	SetSignedCertificate(ctx context.Context, rq *SetSignedCertificateRequest) (*SetSignedCertificateResponse, error)
	GetCertificate(ctx context.Context, serialNumber string) (*CertificateResponse, error)
	Sign(ctx context.Context, roleName string, rq *SignRequest) (*Certificate, error)
	Issue(ctx context.Context, roleName string, rq *IssueRequest) (*IssueResponse, error)
	GetChain(ctx context.Context) (string, error)
	// role CRUD block

	CreateRole(ctx context.Context, rq *RoleRequest) (*Role, error)

	// engines

	UnmountPKIEngine(ctx context.Context) error
	EnablePKIEngine(ctx context.Context) error
}

type vaultImpl struct {
	client *lib.Client
	path   string
	logger kit.CLoggerFunc
}

type Config struct {
	Mock       bool
	Url        string
	Token      string
	TimeoutSec int `config:"timeout-sec"`
}

func (a *vaultImpl) Path(path string) Vault {
	a.path = path
	return a
}

func (a *vaultImpl) Init(ctx context.Context, cfg *Config) error {
	client, err := lib.NewClient(&lib.Config{
		Address: cfg.Url,
		Timeout: time.Second * time.Duration(cfg.TimeoutSec),
	})
	if err != nil {
		return ErrVaultNewClient(ctx, err)
	}
	// authenticate
	client.SetToken(cfg.Token)
	a.client = client
	return nil
}

func NewVault(logger kit.CLoggerFunc) Vault {
	return &vaultImpl{
		logger: logger,
		path:   PathPki,
	}
}

func (a *vaultImpl) l() kit.CLogger {
	return a.logger().Cmp("vault")
}

func (a *vaultImpl) Close(ctx context.Context) error {
	return nil
}

func (a *vaultImpl) GenerateCSR(ctx context.Context, rq *CSRRequest) (*CSR, error) {
	rs := &CSR{}
	if err := a.write(ctx, "intermediate/generate/"+string(rq.Type), &rq.CertificateRequest, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (a *vaultImpl) GetCertificate(ctx context.Context, serialNumber string) (*CertificateResponse, error) {
	rs := &CertificateResponse{}
	if err := a.read(ctx, fmt.Sprintf("cert/%s", serialNumber), &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (a *vaultImpl) GenerateRootCA(ctx context.Context, rq *CertificateRequest) (*Certificate, error) {
	rs := &Certificate{}
	if err := a.write(ctx, "root/generate/internal", rq, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (a *vaultImpl) SignIntermediateCA(ctx context.Context, rq *SignRequest) (*Certificate, error) {
	rs := &Certificate{}
	if err := a.write(ctx, "root/sign-intermediate", rq, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (a *vaultImpl) SetSignedCertificate(ctx context.Context, rq *SetSignedCertificateRequest) (*SetSignedCertificateResponse, error) {
	rs := &SetSignedCertificateResponse{}
	if err := a.write(ctx, "intermediate/set-signed", rq, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (a *vaultImpl) GetChain(ctx context.Context) (string, error) {
	rs := &CertificateResponse{}
	if err := a.read(ctx, "cert/ca_chain", &rs); err != nil {
		return "", err
	}
	return rs.Certificate, nil
}

func (a *vaultImpl) CreateRole(ctx context.Context, rq *RoleRequest) (*Role, error) {
	rs := &Role{}
	if err := a.write(ctx, "roles/"+string(rq.Name), rq, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (a *vaultImpl) Sign(ctx context.Context, roleName string, rq *SignRequest) (*Certificate, error) {
	rs := &Certificate{}
	if err := a.write(ctx, "sign/"+roleName, rq, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (a *vaultImpl) Issue(ctx context.Context, roleName string, rq *IssueRequest) (*IssueResponse, error) {
	rs := &IssueResponse{}
	if err := a.write(ctx, "issue/"+roleName, rq, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (a *vaultImpl) UnmountPKIEngine(ctx context.Context) error {
	err := a.client.Sys().Unmount(PathPki)
	if err != nil {
		return ErrVaultUnmountPKIEngine(ctx, err)
	}
	return nil
}

func (a *vaultImpl) EnablePKIEngine(ctx context.Context) error {
	return a.client.Sys().Mount(a.path, &lib.MountInput{
		Type: EnginePki,
		Config: lib.MountConfigInput{
			MaxLeaseTTL:     MaxTTL,
			DefaultLeaseTTL: TTL,
		},
	})
}

func (a *vaultImpl) toPath(postfix string) string {
	if postfix == "" {
		return a.path
	}
	return fmt.Sprintf("%s/%s", a.path, postfix)
}

func (a *vaultImpl) convertFromMap(ctx context.Context, data, out any) error {
	payload, err := kit.JsonEncode(data)
	if err != nil {
		return err
	}
	if len(payload) == 0 {
		return nil
	}
	err = kit.Unmarshal(payload, &out)
	if err != nil {
		return ErrVaultUnmarshalResponse(ctx, err)
	}
	return nil
}

func (a *vaultImpl) write(ctx context.Context, postfix string, rq any, out any) error {
	data, err := kit.ConvertToMap(rq)
	if err != nil {
		return err
	}
	rs, err := a.client.Logical().Write(a.toPath(postfix), data)
	if err != nil {
		return ErrVaultClientWrite(ctx, err)
	}
	if rs == nil || rs.Data == nil {
		return ErrVaultClientEmptyResponseData(ctx)
	}
	return a.convertFromMap(ctx, rs.Data, &out)
}

func (a *vaultImpl) read(ctx context.Context, postfix string, out any) error {
	rs, err := a.client.Logical().Read(a.toPath(postfix))
	if err != nil {
		return ErrVaultClientRead(ctx, err)
	}
	if rs == nil || rs.Data == nil {
		return ErrVaultClientEmptyResponseData(ctx)
	}
	return a.convertFromMap(ctx, rs.Data, &out)
}
