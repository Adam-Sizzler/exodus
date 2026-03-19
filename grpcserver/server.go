package grpcserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"cerberus-node/config"
	"cerberus-node/proto"
	"cerberus-node/server"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func StartGRPCServer(cfg *config.NodeConfig, nodeServer *server.NodeServer) error {
	var opts []grpc.ServerOption
	opts = append(opts,
		grpc.UnaryInterceptor(grpcUnaryRequestLogger(cfg)),
		grpc.StreamInterceptor(grpcStreamRequestLogger(cfg)),
	)

	var tlsConfig *tls.Config
	if cfg.CERBERUS.MTLSConfig != nil {
		cfg.Logger.Debug("Configuring mTLS for gRPC server")
		cert, err := tls.X509KeyPair(
			[]byte(cfg.CERBERUS.MTLSConfig.Cert),
			[]byte(cfg.CERBERUS.MTLSConfig.Key),
		)
		if err != nil {
			cfg.Logger.Error("Failed to load server certificate", "error", err)
			return err
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM([]byte(cfg.CERBERUS.MTLSConfig.CACert)) {
			cfg.Logger.Error("Failed to parse CA certificate")
			return fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig = &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		cfg.Logger.Info("mTLS enabled for gRPC server")
	} else {
		if cfg.CERBERUS.GrpcAddress != "127.0.0.1" && cfg.CERBERUS.GrpcAddress != "localhost" {
			cfg.Logger.Warn("Insecure gRPC on non-local address", "address", cfg.CERBERUS.GrpcAddress)
		}
		cfg.Logger.Debug("Using insecure gRPC server (no TLS configured inside Node)")
	}

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterNodeServiceServer(grpcServer, nodeServer)

	var pathPrefix string
	if cfg.CERBERUS.GrpcPath != "" {
		pathPrefix = "/" + strings.Trim(cfg.CERBERUS.GrpcPath, "/")
		cfg.Logger.Info("gRPC path prefix configured", "prefix", pathPrefix)
	}

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Logger.Debug("Node HTTP request", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)
		cfg.Logger.Trace(
			"Node HTTP request details",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"x_forwarded_for", r.Header.Get("X-Forwarded-For"),
			"x_forwarded_proto", r.Header.Get("X-Forwarded-Proto"),
		)

		if pathPrefix != "" && strings.HasPrefix(r.URL.Path, pathPrefix) {
			r.URL.Path = strings.TrimPrefix(r.URL.Path, pathPrefix)
		}

		grpcServer.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("%s:%d", cfg.CERBERUS.GrpcAddress, cfg.CERBERUS.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		cfg.Logger.Error("Failed to start listener", "error", err)
		return err
	}

	cfg.Logger.Info("Starting gRPC server", "address", addr)

	httpServer := &http.Server{}
	if tlsConfig != nil {
		httpServer.Handler = mainHandler
		httpServer.TLSConfig = tlsConfig
		if err := http2.ConfigureServer(httpServer, &http2.Server{}); err != nil {
			return fmt.Errorf("configure http2 server: %w", err)
		}
		listener = tls.NewListener(listener, tlsConfig)
	} else {
		httpServer.Handler = h2c.NewHandler(mainHandler, &http2.Server{})
	}

	if err := httpServer.Serve(listener); err != nil {
		cfg.Logger.Error("Failed to serve gRPC", "error", err)
		return err
	}

	return nil
}

func grpcUnaryRequestLogger(cfg *config.NodeConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		remoteAddr := ""
		if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
			remoteAddr = p.Addr.String()
		}

		cfg.Logger.Debug("gRPC unary request", "method", info.FullMethod, "remote_addr", remoteAddr)
		cfg.Logger.Trace("gRPC unary request details", "method", info.FullMethod, "payload_type", fmt.Sprintf("%T", req))

		resp, err := handler(ctx, req)
		durationMs := time.Since(start).Milliseconds()
		code := status.Code(err).String()

		if err != nil {
			cfg.Logger.Warn("gRPC unary failed", "method", info.FullMethod, "code", code, "duration_ms", durationMs, "error", err)
			return resp, err
		}

		cfg.Logger.Debug("gRPC unary completed", "method", info.FullMethod, "code", code, "duration_ms", durationMs)
		return resp, nil
	}
}

func grpcStreamRequestLogger(cfg *config.NodeConfig) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		remoteAddr := ""
		if p, ok := peer.FromContext(ss.Context()); ok && p.Addr != nil {
			remoteAddr = p.Addr.String()
		}

		cfg.Logger.Debug(
			"gRPC stream opened",
			"method", info.FullMethod,
			"client_stream", info.IsClientStream,
			"server_stream", info.IsServerStream,
			"remote_addr", remoteAddr,
		)
		cfg.Logger.Trace("gRPC stream details", "method", info.FullMethod)

		err := handler(srv, ss)
		durationMs := time.Since(start).Milliseconds()
		code := status.Code(err).String()
		if err != nil {
			cfg.Logger.Warn("gRPC stream failed", "method", info.FullMethod, "code", code, "duration_ms", durationMs, "error", err)
			return err
		}

		cfg.Logger.Debug("gRPC stream closed", "method", info.FullMethod, "code", code, "duration_ms", durationMs)
		return nil
	}
}
