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
	"strings"
	"time"

	"github.com/exodus/subscription-page/backend/internal/config"
	"github.com/exodus/subscription-page/backend/internal/logger"
	"github.com/exodus/subscription-page/backend/internal/proto"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const grpcTokenHeader = "x-exodus-grpc-token"

func Start(ctx context.Context, cfg config.Config, nodeService proto.NodeServiceServer) error {
	log := logger.WithContext("GrpcServer")
	if nodeService == nil {
		return fmt.Errorf("node service is required")
	}

	expectedToken := strings.TrimSpace(cfg.GRPCToken)
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpcPathValidationUnaryInterceptor(log),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		grpcPathValidationStreamInterceptor(log),
	}
	opts := make([]grpc.ServerOption, 0, 2)
	if cfg.RequireGRPCToken {
		if expectedToken == "" {
			return config.NewEnvError("SUB_GRPC_TOKEN", "Required when SUB_SECRET_KEY is not provided. Dashboard → Subscription → Nodes → Current node → gRPC Token (SUB_GRPC_TOKEN) or Secret Key (SUB_SECRET_KEY).")
		}
		unaryInterceptors = append(unaryInterceptors, grpcTokenUnaryInterceptor(expectedToken, log))
		streamInterceptors = append(streamInterceptors, grpcTokenStreamInterceptor(expectedToken, log))
		log.Info("[CONFIG] gRPC auth mode: TLS + token")
	} else {
		log.Info("[CONFIG] gRPC auth mode: mTLS")
	}
	opts = append(opts,
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterNodeServiceServer(grpcServer, nodeService)

	pathPrefix := cfg.SubPathTrimmed()
	if pathPrefix != "" {
		log.Info("[CONFIG] gRPC path prefix: " + pathPrefix)
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

	httpServer := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}
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
		log.Info("[CONFIG] gRPC mTLS enabled")
	} else {
		httpServer.Handler = h2c.NewHandler(handler, &http2.Server{})
		log.Info("[CONFIG] gRPC h2c mode enabled (for reverse proxy TLS termination)")
	}

	log.Info("[CONFIG] gRPC listening on " + addr)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warn("gRPC shutdown failed", logger.String("error", err.Error()))
		}
	}()

	err = httpServer.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
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

func grpcPathValidationUnaryInterceptor(log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if info == nil || len(info.FullMethod) == 0 || info.FullMethod[0] != '/' {
			log.Warn("Malformed gRPC unary method name")
			return nil, status.Error(codes.Unimplemented, "malformed method name")
		}
		return handler(ctx, req)
	}
}

func grpcPathValidationStreamInterceptor(log *logger.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if info == nil || len(info.FullMethod) == 0 || info.FullMethod[0] != '/' {
			log.Warn("Malformed gRPC stream method name")
			return status.Error(codes.Unimplemented, "malformed method name")
		}
		return handler(srv, ss)
	}
}

func grpcTokenUnaryInterceptor(expectedToken string, log *logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := validateIncomingGRPCToken(ctx, expectedToken); err != nil {
			method := "unknown"
			if info != nil {
				method = info.FullMethod
			}
			log.Warn("gRPC unary authentication failed", logger.String("method", method), logger.String("error", err.Error()))
			return nil, err
		}
		return handler(ctx, req)
	}
}

func grpcTokenStreamInterceptor(expectedToken string, log *logger.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := validateIncomingGRPCToken(ss.Context(), expectedToken); err != nil {
			method := "unknown"
			if info != nil {
				method = info.FullMethod
			}
			log.Warn("gRPC stream authentication failed", logger.String("method", method), logger.String("error", err.Error()))
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
