package users

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	dbmanager "exodus/backend/db/manager"
)

func (nm *NodeMonitor) loadNodeMTLSConfig(ctx context.Context) (*tls.Config, error) {
	var (
		caCertPEM     string
		clientCertPEM string
		clientKeyPEM  string
	)

	err := nm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT ca_cert, client_cert, client_key
			FROM keygen
			ORDER BY created_at ASC
			LIMIT 1
		`).Scan(&caCertPEM, &clientCertPEM, &clientKeyPEM)
	})
	if err != nil {
		return nil, fmt.Errorf("load keygen mTLS material: %w", err)
	}

	clientCert, err := tls.X509KeyPair([]byte(clientCertPEM), []byte(clientKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("parse keygen client certificate/key: %w", err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(caCertPEM)) {
		return nil, fmt.Errorf("parse keygen CA certificate")
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		RootCAs:      pool,
		Certificates: []tls.Certificate{clientCert},
		ServerName:   "internal.exodus.local",
	}, nil
}
