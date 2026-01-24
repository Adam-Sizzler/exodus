package grpcserver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"v2ray-stat/node/config"
	"v2ray-stat/node/server"
	"v2ray-stat/proto"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func StartGRPCServer(cfg *config.NodeConfig, nodeServer *server.NodeServer) error {
	var opts []grpc.ServerOption

	if cfg.V2RS.MTLSConfig != nil {
		cfg.Logger.Debug("Configuring mTLS for gRPC server")
		cert, err := tls.LoadX509KeyPair(cfg.V2RS.MTLSConfig.Cert, cfg.V2RS.MTLSConfig.Key)
		if err != nil {
			cfg.Logger.Error("Failed to load server certificate", "error", err)
			return err
		}
		caCert, err := os.ReadFile(cfg.V2RS.MTLSConfig.CACert)
		if err != nil {
			cfg.Logger.Error("Failed to read CA certificate", "error", err)
			return err
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			cfg.Logger.Error("Failed to parse CA certificate")
			return fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.Creds(creds))
		cfg.Logger.Info("mTLS enabled for gRPC server")
	} else {
		if cfg.V2RS.GrpcAddress != "127.0.0.1" && cfg.V2RS.GrpcAddress != "localhost" {
			cfg.Logger.Warn("Insecure gRPC on non-local address", "address", cfg.V2RS.GrpcAddress)
		}
		cfg.Logger.Debug("Using insecure gRPC server (no TLS configured inside Node)")
	}

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterNodeServiceServer(grpcServer, nodeServer)

	var pathPrefix string
	if cfg.V2RS.GrpcPath != "" {
		pathPrefix = "/" + strings.Trim(cfg.V2RS.GrpcPath, "/")
		cfg.Logger.Info("gRPC path prefix configured", "prefix", pathPrefix)
	}

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pathPrefix != "" && strings.HasPrefix(r.URL.Path, pathPrefix) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, pathPrefix)
		}

		grpcServer.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("%s:%d", cfg.V2RS.GrpcAddress, cfg.V2RS.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		cfg.Logger.Error("Failed to start listener", "error", err)
		return err
	}

	cfg.Logger.Info("Starting gRPC server", "address", addr)

	h2Server := &http2.Server{}
	httpServer := &http.Server{
		Handler: h2c.NewHandler(mainHandler, h2Server),
	}

	if err := httpServer.Serve(listener); err != nil {
		cfg.Logger.Error("Failed to serve gRPC", "error", err)
		return err
	}

	return nil
}
