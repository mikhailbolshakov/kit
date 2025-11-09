package mongodb

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"gitlab.com/algmib/kit"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"strings"
	"time"
)

type Config struct {
	ConnectionString string
	TimeoutSec       *int
	CertPath         *string
}

type Storage struct {
	Instance *mongo.Client
	lg       kit.CLoggerFunc
}

func Open(config *Config, logger kit.CLoggerFunc) (*Storage, error) {

	s := &Storage{lg: logger}

	l := logger().Cmp("mongo").Mth("open").Dbg("connecting...")

	ctx := context.Background()

	// setup connection timeout if specified
	if config.TimeoutSec != nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(*config.TimeoutSec)*time.Second)
		defer cancel()
	}

	// setup options
	opts := options.Client()
	opts = opts.ApplyURI(config.ConnectionString)

	// tls configuration
	tlsConfig, err := s.makeTlsConfig(config.CertPath)
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		opts.SetTLSConfig(tlsConfig)
	}

	// connect
	s.Instance, err = mongo.Connect(ctx, opts)
	if err != nil {
		return nil, ErrConnection(err)
	}

	l.Dbg("ok")

	return s, nil
}

func (s *Storage) Close(ctx context.Context) {
	_ = s.Instance.Disconnect(ctx)
}

func (s *Storage) makeTlsConfig(dbCertPath *string) (*tls.Config, error) {

	if dbCertPath == nil || strings.Compare(*dbCertPath, "") == 0 {
		return nil, nil
	}

	rootPEM, err := os.ReadFile(*dbCertPath)
	if err != nil {
		return nil, ErrReadCertFile(err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootPEM)
	if !ok {
		return nil, ErrAppendCert(err)
	}

	return &tls.Config{
		RootCAs:            roots,
		InsecureSkipVerify: true,
	}, nil
}
