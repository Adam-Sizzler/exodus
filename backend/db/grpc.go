package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"v2ray-stat/backend/config"
	"v2ray-stat/proto"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// NodeClient represents a gRPC client for a node.
type NodeClient struct {
	NodeName string
	URL      string
	Client   proto.NodeServiceClient
	Conn     *grpc.ClientConn
}

func pathPrefixUnaryInterceptor(prefix string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newMethod := prefix + method
		return invoker(ctx, newMethod, req, reply, cc, opts...)
	}
}

func pathPrefixStreamInterceptor(prefix string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		newMethod := prefix + method
		return streamer(ctx, desc, cc, newMethod, opts...)
	}
}

func NewNodeClient(nodeCfg config.NodeConfig, cfg *config.BackendConfig) (*NodeClient, error) {
	url := fmt.Sprintf("%s:%s", nodeCfg.Address, nodeCfg.Port)
	var opts []grpc.DialOption

	if nodeCfg.Path != "" {
		cleanPath := "/" + strings.Trim(nodeCfg.Path, "/")
		opts = append(opts, grpc.WithUnaryInterceptor(pathPrefixUnaryInterceptor(cleanPath)))
		opts = append(opts, grpc.WithStreamInterceptor(pathPrefixStreamInterceptor(cleanPath)))
		cfg.Logger.Debug("Using gRPC path prefix", "node", nodeCfg.NodeName, "prefix", cleanPath)
	}

	useTLS := nodeCfg.Schema == "https" || nodeCfg.MTLSConfig != nil

	if useTLS {
		tlsConfig := &tls.Config{
			ServerName: nodeCfg.Address, // Важно для SNI (Nginx требует правильный хост)
		}

		if nodeCfg.MTLSConfig != nil {
			if nodeCfg.MTLSConfig.CACert != "" {
				caCert, err := os.ReadFile(nodeCfg.MTLSConfig.CACert)
				if err != nil {
					return nil, fmt.Errorf("failed to read CA cert: %v", err)
				}
				certPool := x509.NewCertPool()
				certPool.AppendCertsFromPEM(caCert)
				tlsConfig.RootCAs = certPool
			} else {
				tlsConfig.InsecureSkipVerify = true
			}

			if nodeCfg.MTLSConfig.Cert != "" && nodeCfg.MTLSConfig.Key != "" {
				cert, err := tls.LoadX509KeyPair(nodeCfg.MTLSConfig.Cert, nodeCfg.MTLSConfig.Key)
				if err != nil {
					return nil, fmt.Errorf("failed to load client cert: %v", err)
				}
				tlsConfig.Certificates = []tls.Certificate{cert}
			}
		} else {
			// Обычный HTTPS (Public Trusted CA, например Let's Encrypt)
			// Оставляем RootCAs nil, Go использует системные
		}

		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Создаем подключение
	// grpc.NewClient
	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node %s: %v", nodeCfg.NodeName, err)
	}

	return &NodeClient{
		NodeName: nodeCfg.NodeName,
		URL:      url,
		Client:   proto.NewNodeServiceClient(conn),
		Conn:     conn,
	}, nil
}

// Метод для закрытия соединения
func (nc *NodeClient) Close() error {
	if nc.Conn != nil {
		return nc.Conn.Close()
	}
	return nil
}

// InitNodeClients initializes gRPC clients for all nodes.
func InitNodeClients(cfg *config.BackendConfig) ([]*NodeClient, error) {
	var nodeClients []*NodeClient
	for _, node := range cfg.Nodes {
		client, err := NewNodeClient(node, cfg)
		if err != nil {
			cfg.Logger.Error("Failed to initialize client for node", "node", node.NodeName, "error", err)
			continue
		}
		nodeClients = append(nodeClients, client)
	}

	if len(nodeClients) == 0 {
		return nil, fmt.Errorf("no node clients initialized")
	}
	return nodeClients, nil
}
