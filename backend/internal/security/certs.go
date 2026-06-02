package security

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type MasterCerts struct {
	CACertPEM     string
	CAKeyPEM      string
	ClientCertPEM string
	ClientKeyPEM  string
}

type NodeCert struct {
	NodeCertPEM string
	NodeKeyPEM  string
}

func GenerateGRPCAuthToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate grpc auth token: %w", err)
	}
	return hex.EncodeToString(raw), nil
}

func ResolveGRPCAuthToken(value string) (string, error) {
	token := strings.ToLower(strings.TrimSpace(value))
	if token == "" {
		return GenerateGRPCAuthToken()
	}

	raw, err := hex.DecodeString(token)
	if err != nil || len(raw) != 32 {
		return "", fmt.Errorf("grpc auth token must be a 64-character hexadecimal string")
	}
	return token, nil
}

func GenerateJWTKeypair() (publicKeyPEM string, privateKeyPEM string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("generate rsa private key: %w", err)
	}

	privateDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", "", fmt.Errorf("marshal rsa private key: %w", err)
	}
	publicDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("marshal rsa public key: %w", err)
	}

	privatePEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateDER})
	publicPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER})
	return string(publicPEM), string(privatePEM), nil
}

func GenerateMasterCerts() (MasterCerts, error) {
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return MasterCerts{}, fmt.Errorf("generate ca key: %w", err)
	}

	caTemplate, err := newCertificateTemplate()
	if err != nil {
		return MasterCerts{}, err
	}
	caTemplate.Subject = pkix.Name{CommonName: randomCertCN()}
	caTemplate.IsCA = true
	caTemplate.BasicConstraintsValid = true
	caTemplate.KeyUsage = x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature
	caTemplate.ExtKeyUsage = nil

	caDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return MasterCerts{}, fmt.Errorf("create ca certificate: %w", err)
	}

	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return MasterCerts{}, fmt.Errorf("generate client key: %w", err)
	}
	clientTemplate, err := newCertificateTemplate()
	if err != nil {
		return MasterCerts{}, err
	}
	clientTemplate.Subject = pkix.Name{CommonName: randomCertCN()}
	clientTemplate.KeyUsage = x509.KeyUsageDigitalSignature
	clientTemplate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	clientDER, err := x509.CreateCertificate(rand.Reader, clientTemplate, caTemplate, &clientKey.PublicKey, caKey)
	if err != nil {
		return MasterCerts{}, fmt.Errorf("create client certificate: %w", err)
	}

	caKeyDER, err := x509.MarshalPKCS8PrivateKey(caKey)
	if err != nil {
		return MasterCerts{}, fmt.Errorf("marshal ca private key: %w", err)
	}
	clientKeyDER, err := x509.MarshalPKCS8PrivateKey(clientKey)
	if err != nil {
		return MasterCerts{}, fmt.Errorf("marshal client private key: %w", err)
	}

	return MasterCerts{
		CACertPEM:     string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})),
		CAKeyPEM:      string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: caKeyDER})),
		ClientCertPEM: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientDER})),
		ClientKeyPEM:  string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: clientKeyDER})),
	}, nil
}

func GenerateNodeCert(caCertPEM string, caKeyPEM string) (NodeCert, error) {
	caCertBlock, _ := pem.Decode([]byte(caCertPEM))
	if caCertBlock == nil {
		return NodeCert{}, fmt.Errorf("invalid ca certificate pem")
	}
	caCert, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		return NodeCert{}, fmt.Errorf("parse ca certificate: %w", err)
	}

	caKeyBlock, _ := pem.Decode([]byte(caKeyPEM))
	if caKeyBlock == nil {
		return NodeCert{}, fmt.Errorf("invalid ca key pem")
	}
	caKeyRaw, err := x509.ParsePKCS8PrivateKey(caKeyBlock.Bytes)
	if err != nil {
		return NodeCert{}, fmt.Errorf("parse ca key pkcs8: %w", err)
	}
	caKey, ok := caKeyRaw.(*ecdsa.PrivateKey)
	if !ok {
		return NodeCert{}, fmt.Errorf("unsupported ca private key type")
	}

	nodeKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return NodeCert{}, fmt.Errorf("generate node key: %w", err)
	}
	nodeTemplate, err := newCertificateTemplate()
	if err != nil {
		return NodeCert{}, err
	}
	nodeTemplate.Subject = pkix.Name{CommonName: randomCertCN()}
	nodeTemplate.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	nodeTemplate.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	nodeTemplate.DNSNames = []string{"internal.exodus.local"}

	nodeDER, err := x509.CreateCertificate(rand.Reader, nodeTemplate, caCert, &nodeKey.PublicKey, caKey)
	if err != nil {
		return NodeCert{}, fmt.Errorf("create node certificate: %w", err)
	}

	nodeKeyDER, err := x509.MarshalPKCS8PrivateKey(nodeKey)
	if err != nil {
		return NodeCert{}, fmt.Errorf("marshal node private key: %w", err)
	}

	return NodeCert{
		NodeCertPEM: string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: nodeDER})),
		NodeKeyPEM:  string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: nodeKeyDER})),
	}, nil
}

func newCertificateTemplate() (*x509.Certificate, error) {
	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return nil, fmt.Errorf("generate certificate serial: %w", err)
	}
	now := time.Now().UTC()
	return &x509.Certificate{
		SerialNumber: serial,
		NotBefore:    now.Add(-5 * time.Minute),
		NotAfter:     now.AddDate(3, 0, 0),
	}, nil
}

func randomCertCN() string {
	alphabet := []byte("0123456789ABCDEFGHJKLMNPQRSTUVWXYZ_abcdefghjkmnopqrstuvwxyz-")
	size := 24
	out := make([]byte, size)
	for i := range out {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			out[i] = alphabet[0]
			continue
		}
		out[i] = alphabet[n.Int64()]
	}
	return string(out)
}
