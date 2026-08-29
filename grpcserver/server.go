package grpcserver

import (
	"context"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"exodus-node/config"
	"exodus-node/proto"
	"exodus-node/server"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const grpcTokenHeader = "x-exodus-grpc-token"

func StartGRPCServer(cfg *config.NodeConfig, nodeServer *server.NodeServer) error {
	if cfg == nil {
		return fmt.Errorf("node config is nil")
	}
	if cfg.Logger == nil {
		return fmt.Errorf("node logger is nil")
	}
	if nodeServer == nil {
		return fmt.Errorf("node server is nil")
	}
	log := cfg.LoggerFor("GrpcServer")

	expectedToken := strings.TrimSpace(cfg.Backend.GRPCToken)
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpcPathValidationUnaryInterceptor(),
		grpcUnaryRequestLogger(log),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		grpcPathValidationStreamInterceptor(),
		grpcStreamRequestLogger(log),
	}
	if cfg.Backend.RequireGRPCToken {
		if expectedToken == "" {
			return fmt.Errorf("NODE_GRPC_TOKEN is required")
		}
		unaryInterceptors = append(unaryInterceptors, grpcTokenUnaryInterceptor(expectedToken))
		streamInterceptors = append(streamInterceptors, grpcTokenStreamInterceptor(expectedToken))
		log.Debug("gRPC auth mode: TLS + token")
	} else {
		log.Debug("gRPC auth mode: mTLS")
	}

	var opts []grpc.ServerOption
	opts = append(opts,
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	var tlsConfig *tls.Config
	if cfg.Backend.MTLSConfig != nil {
		log.Debug("Configuring mTLS for gRPC server")
		cert, err := tls.X509KeyPair(
			[]byte(cfg.Backend.MTLSConfig.Cert),
			[]byte(cfg.Backend.MTLSConfig.Key),
		)
		if err != nil {
			log.Error("Failed to load server certificate", "error", err)
			return err
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM([]byte(cfg.Backend.MTLSConfig.CACert)) {
			log.Error("Failed to parse CA certificate")
			return fmt.Errorf("failed to parse CA certificate")
		}
		tlsConfig = &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}
		if cfg.Backend.SNIVerification {
			expectedSNI := ""
			if cfg.SecretKeyReport != nil {
				expectedSNI = cfg.SecretKeyReport.SNI
			}
			log.Info("Strict SNI verification enabled for mTLS", "expected_sni", expectedSNI)
			tlsConfig.GetConfigForClient = func(hello *tls.ClientHelloInfo) (*tls.Config, error) {
				if expectedSNI != "" && !strings.EqualFold(strings.TrimSpace(hello.ServerName), expectedSNI) {
					log.Warn("Rejected mTLS connection: invalid SNI", "received_sni", hello.ServerName, "expected_sni", expectedSNI)
					return nil, fmt.Errorf("tls: invalid server name %q, expected %q", hello.ServerName, expectedSNI)
				}
				return nil, nil
			}
		} else {
			log.Debug("mTLS SNI verification disabled (standard mTLS mode)")
		}
		log.Debug("mTLS enabled for gRPC server")
	} else {
		if cfg.Backend.GrpcAddress != "127.0.0.1" && cfg.Backend.GrpcAddress != "localhost" {
			log.Warn("Insecure gRPC on non-local address", "address", cfg.Backend.GrpcAddress)
		}
		log.Debug("Using h2c gRPC server (for reverse proxy TLS termination)")
	}

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterNodeServiceServer(grpcServer, nodeServer)

	pathPrefix := cfg.Backend.Trimmed()
	if pathPrefix != "" {
		log.Debug("gRPC path prefix configured", "prefix", pathPrefix)
	}

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Node HTTP request", "method", r.Method, "path", r.URL.Path, "remote_addr", r.RemoteAddr)
		log.Trace(
			"Node HTTP request details",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"remote_addr", r.RemoteAddr,
			"user_agent", r.UserAgent(),
			"x_forwarded_for", r.Header.Get("X-Forwarded-For"),
			"x_forwarded_proto", r.Header.Get("X-Forwarded-Proto"),
		)

		if pathPrefix != "" {
			var ok bool
			r.URL.Path, ok = trimPathPrefix(r.URL.Path, pathPrefix)
			if !ok {
				http.NotFound(w, r)
				return
			}
		}

		grpcServer.ServeHTTP(w, r)
	})

	addr := fmt.Sprintf("%s:%d", cfg.Backend.GrpcAddress, cfg.Backend.GrpcPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Failed to start listener", "error", err)
		return err
	}

	log.Debug("Starting gRPC server", "address", addr)

	httpServer := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info("Shutting down gRPC server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
	}()

	if err := httpServer.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("Failed to serve gRPC", "error", err)
		return err
	}

	return nil
}

func trimPathPrefix(path string, prefix string) (string, bool) {
	if prefix == "" {
		return path, true
	}
	if path == prefix {
		return "/", true
	}
	if strings.HasPrefix(path, prefix+"/") {
		return strings.TrimPrefix(path, prefix), true
	}
	return path, false
}

func grpcPathValidationUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if info == nil || len(info.FullMethod) == 0 || info.FullMethod[0] != '/' {
			return nil, status.Error(codes.Unimplemented, "malformed method name")
		}
		return handler(ctx, req)
	}
}

func grpcPathValidationStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if info == nil || len(info.FullMethod) == 0 || info.FullMethod[0] != '/' {
			return status.Error(codes.Unimplemented, "malformed method name")
		}
		return handler(srv, ss)
	}
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

func grpcUnaryRequestLogger(log *config.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		remoteAddr := server.ExtractClientIP(ctx)

		log.Debug("gRPC unary request", "method", info.FullMethod, "remote_addr", remoteAddr)
		log.Trace("gRPC unary request details", "method", info.FullMethod, "payload_type", fmt.Sprintf("%T", req))

		resp, err := handler(ctx, req)
		durationMs := time.Since(start).Milliseconds()
		code := status.Code(err).String()

		if err != nil {
			log.Warn("gRPC unary failed", "method", info.FullMethod, "code", code, "duration_ms", durationMs, "error", err)
			return resp, err
		}

		log.Debug("gRPC unary completed", "method", info.FullMethod, "code", code, "duration_ms", durationMs)
		return resp, nil
	}
}

func grpcStreamRequestLogger(log *config.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		remoteAddr := server.ExtractClientIP(ss.Context())

		log.Debug(
			"gRPC stream opened",
			"method", info.FullMethod,
			"client_stream", info.IsClientStream,
			"server_stream", info.IsServerStream,
			"remote_addr", remoteAddr,
		)
		log.Trace("gRPC stream details", "method", info.FullMethod)

		err := handler(srv, ss)
		durationMs := time.Since(start).Milliseconds()
		code := status.Code(err).String()
		if err != nil {
			log.Warn("gRPC stream failed", "method", info.FullMethod, "code", code, "duration_ms", durationMs, "error", err)
			return err
		}

		log.Debug("gRPC stream closed", "method", info.FullMethod, "code", code, "duration_ms", durationMs)
		return nil
	}
}
