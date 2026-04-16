package grpcserver

import (
	"context"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/exodus/subscription-page/backend/internal/config"
	"github.com/exodus/subscription-page/backend/internal/proto"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const grpcTokenHeader = "x-exodus-grpc-token"

func Start(cfg config.Config, nodeService proto.NodeServiceServer) error {
	expectedToken := strings.TrimSpace(cfg.GRPCToken)
	opts := make([]grpc.ServerOption, 0, 2)
	if cfg.RequireGRPCToken {
		if expectedToken == "" {
			return fmt.Errorf("SUB_GRPC_TOKEN is required")
		}
		opts = append(opts,
			grpc.UnaryInterceptor(grpcTokenUnaryInterceptor(expectedToken)),
			grpc.StreamInterceptor(grpcTokenStreamInterceptor(expectedToken)),
		)
		log.Printf("[CONFIG] gRPC auth mode: TLS + token")
	} else {
		log.Printf("[CONFIG] gRPC auth mode: mTLS")
	}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterNodeServiceServer(grpcServer, nodeService)

	pathPrefix := "/" + strings.Trim(strings.TrimSpace(cfg.GRPCPath), "/")
	if pathPrefix == "/" {
		pathPrefix = ""
	} else {
		log.Printf("[CONFIG] gRPC path prefix: %s", pathPrefix)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pathPrefix != "" {
			trimmedPath, ok := trimGRPCPathPrefix(r.URL.Path, pathPrefix)
			if !ok {
				abortConnection(w)
				return
			}
			r.URL.Path = trimmedPath
		}
		grpcServer.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("%s:%d", cfg.GRPCAddress, cfg.GRPCPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen grpc: %w", err)
	}

	httpServer := &http.Server{}
	if cfg.MTLSConfig != nil {
		tlsCfg, tlsErr := buildServerTLSConfig(cfg.MTLSConfig)
		if tlsErr != nil {
			return tlsErr
		}
		httpServer.Handler = handler
		httpServer.TLSConfig = tlsCfg
		if err := http2.ConfigureServer(httpServer, &http2.Server{}); err != nil {
			return fmt.Errorf("configure http2: %w", err)
		}
		listener = tls.NewListener(listener, tlsCfg)
		log.Printf("[CONFIG] gRPC mTLS enabled")
	} else {
		httpServer.Handler = h2c.NewHandler(handler, &http2.Server{})
		log.Printf("[CONFIG] gRPC h2c mode enabled (for reverse proxy TLS termination)")
	}

	log.Printf("[CONFIG] gRPC listening on %s", addr)
	return httpServer.Serve(listener)
}

func buildServerTLSConfig(material *config.MTLSConfig) (*tls.Config, error) {
	cert, err := tls.X509KeyPair([]byte(material.Cert), []byte(material.Key))
	if err != nil {
		return nil, fmt.Errorf("parse server certificate/key: %w", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM([]byte(material.CACert)) {
		return nil, fmt.Errorf("parse CA certificate")
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil
}

func grpcTokenUnaryInterceptor(expectedToken string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := validateIncomingGRPCToken(ctx, expectedToken); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func grpcTokenStreamInterceptor(expectedToken string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := validateIncomingGRPCToken(ss.Context(), expectedToken); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func validateIncomingGRPCToken(ctx context.Context, expectedToken string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing grpc metadata")
	}

	values := md.Get(grpcTokenHeader)
	if len(values) == 0 {
		return status.Error(codes.Unauthenticated, "missing grpc token")
	}

	providedToken := strings.TrimSpace(values[0])
	if providedToken == "" {
		return status.Error(codes.Unauthenticated, "missing grpc token")
	}
	if subtle.ConstantTimeCompare([]byte(providedToken), []byte(expectedToken)) != 1 {
		return status.Error(codes.Unauthenticated, "invalid grpc token")
	}

	return nil
}

func trimGRPCPathPrefix(path, prefix string) (string, bool) {
	if prefix == "" {
		if path == "" {
			return "/", true
		}
		return path, true
	}

	if path == prefix {
		return "/", true
	}
	withSlash := prefix + "/"
	if !strings.HasPrefix(path, withSlash) {
		return "", false
	}

	trimmed := strings.TrimPrefix(path, prefix)
	if trimmed == "" {
		return "/", true
	}
	return trimmed, true
}

func abortConnection(w http.ResponseWriter) {
	if hijacker, ok := w.(http.Hijacker); ok {
		conn, _, err := hijacker.Hijack()
		if err == nil {
			_ = conn.Close()
			return
		}
	}
	panic(http.ErrAbortHandler)
}
