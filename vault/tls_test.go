//go:build integration

package vault

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/mikhailbolshakov/kit"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type ocppSuite struct {
	kit.Suite
}

func (s *ocppSuite) SetupSuite() {
	s.Suite.Init(func() kit.CLogger { return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel})) })
}

func TestOCPPSuite(t *testing.T) {
	suite.Run(t, new(ocppSuite))
}

func (s *ocppSuite) SetupTest() {
}

func (s *ocppSuite) Test_ClientServerSideTlsSecure() {
	client := s.create(s.Ctx, "pki", "http://localhost:8201", "00000000-0000-0000-0000-000000000000")
	server := client //create(ctx, "pki", "http://localhost:8202", "00000000-0000-0000-0000-000000000000")

	root, err := server.GenerateRootCA(s.Ctx, &CertificateRequest{CommonName: "localhost"})
	s.NoError(err)

	// generate CSR and get certificate
	csrClient, err := client.GenerateCSR(s.Ctx, &CSRRequest{
		Type:               CSRTypeExported,
		CertificateRequest: CertificateRequest{CommonName: "localhost"},
	})
	s.NoError(err)

	csrServer, err := client.GenerateCSR(s.Ctx, &CSRRequest{
		Type:               CSRTypeExported,
		CertificateRequest: CertificateRequest{CommonName: "localhost"},
	})
	s.NoError(err)

	certClient, err := server.SignIntermediateCA(s.Ctx, &SignRequest{
		CSR:      csrClient.CSR,
		Format:   SignFormatPemBundle,
		NotAfter: ToDateFormat(kit.Now().Add(time.Hour)),
	})
	s.NoError(err)

	certServer, err := server.SignIntermediateCA(s.Ctx, &SignRequest{
		CSR:      csrServer.CSR,
		Format:   SignFormatPemBundle,
		NotAfter: ToDateFormat(kit.Now().Add(time.Hour)),
	})
	s.NoError(err)

	cl, err := ParseCertificate(s.Ctx, certClient.Certificate)
	sr, err := ParseCertificate(s.Ctx, certServer.Certificate)
	rt, err := ParseCertificate(s.Ctx, root.Certificate)

	s.NoError(VerifyCertificate(s.Ctx, cl, rt))
	s.NoError(VerifyCertificate(s.Ctx, sr, rt))
	// now client have a certificate and can use it to sing to the server
	// TODO: put root to ca, to be able to validate server certificate

	tlsServerCfg, err := createServerTlsConfig(root.Certificate, certServer.Certificate, csrServer.PrivateKey)
	s.NoError(err)

	tlsClientCfg, err := createClientTlsConfig(root.Certificate, certClient.Certificate, csrClient.PrivateKey)
	s.NoError(err)

	httpServer, err := s.serverHttp(":8443", tlsServerCfg)
	s.NoError(err)
	defer httpServer.Close()
	wsServer, err := s.serverWebSocket(":8444", tlsServerCfg)
	s.NoError(err)
	defer wsServer.Close()

	time.Sleep(time.Second)

	s.NoError(clientWebSocket(":8444", tlsClientCfg))
	s.NoError(clientHttp(":8443", tlsClientCfg))

	time.Sleep(time.Second)
}

func createServerTlsConfig(root, certificate, privateKey string) (*tls.Config, error) {
	rootCertPool, cert, err := parseCerts(root, certificate, privateKey)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{*cert},
		ClientCAs:    rootCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert, // Require client certificates
		MinVersion:   tls.VersionTLS12,
	}, nil
}

func parseCerts(root, certificate, privateKey string) (*x509.CertPool, *tls.Certificate, error) {
	cert, err := tls.X509KeyPair([]byte(certificate), []byte(privateKey))
	if err != nil {
		return nil, nil, err
	}

	rootCertPool := x509.NewCertPool()
	ok := rootCertPool.AppendCertsFromPEM([]byte(root))
	if !ok {
		return nil, nil, errors.New("not ok")
	}
	return rootCertPool, &cert, nil
}

func createClientTlsConfig(root, certificate, privateKey string) (*tls.Config, error) {
	rootCertPool, cert, err := parseCerts(root, certificate, privateKey)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{*cert},
		RootCAs:            rootCertPool,
		ClientCAs:          rootCertPool,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: false, // Set this to false to enable server certificate validation
	}, nil
}

func (s *ocppSuite) serverWebSocket(addr string, tlsCfg *tls.Config) (*http.Server, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.wsHandler)
	server := &http.Server{
		Addr:      addr, // Listen on port 8443 for TLS
		Handler:   mux,
		TLSConfig: tlsCfg,
	}
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
			l().E(err).Err()
		}
	}()
	return server, nil
}

func (s *ocppSuite) serverHttp(addr string, tlsCfg *tls.Config) (*http.Server, error) {
	server := &http.Server{
		Addr:      addr,
		TLSConfig: tlsCfg,
	}

	// Handle requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		l().DbgF("[HTTP] Server: request processing")
	})

	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil {
			l().E(err).Err()
		}
	}()
	return server, nil
}

func clientHttp(addr string, tlsCfg *tls.Config) error {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsCfg,
		},
	}

	// Send a request to the server
	resp, err := client.Get("https://localhost" + addr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	l().DbgF("[HTTP] Client: response - %s\n", body)
	return nil
}

func clientWebSocket(addr string, tlsCfg *tls.Config) error {
	dialer := websocket.Dialer{
		TLSClientConfig: tlsCfg,
	}
	conn, _, err := dialer.Dial("wss://localhost"+addr+"/ws", nil) // "wss" is for WebSocket Secure
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send a message to the server
	err = conn.WriteMessage(websocket.TextMessage, []byte("[WS] Client: request"))
	if err != nil {
		return err
	}

	// Read the server's response
	_, message, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	// Print out the server's response
	l().DbgF("[WS] Client: response: %s\n", message)
	return nil
}
func l() kit.CLogger {
	return kit.L(kit.InitLogger(&kit.LogConfig{Level: kit.TraceLevel, Format: kit.FormatterJson})).
		Srv("ws-client").Node(kit.NewRandString())
}

func (s *ocppSuite) create(ctx context.Context, path, url, token string) Vault {
	client := NewVault(l).Path(path)
	err := client.Init(ctx, &Config{
		Url:        url,
		Token:      token,
		TimeoutSec: 10,
	})
	s.NoError(err)

	// prepare environment
	_ = client.UnmountPKIEngine(ctx)
	_ = client.EnablePKIEngine(ctx)

	return client
}

var upgrader = websocket.Upgrader{
	// Allow connections from any origin for now (adjust for production security)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *ocppSuite) wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	s.NoError(err)
	defer conn.Close()

	// Read message from the WebSocket client
	messageType, message, err := conn.ReadMessage()
	s.NoError(err)

	l().DbgF("[WS] Server: request - %s\n", message)

	// Echo the message back to the client
	err = conn.WriteMessage(messageType, message)
	s.NoError(err)
}
