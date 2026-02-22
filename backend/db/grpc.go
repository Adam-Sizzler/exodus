package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
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
	url := fmt.Sprintf("%s:%d", nodeCfg.Address, nodeCfg.Port)
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
			ServerName: nodeCfg.Address,
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
		}

		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

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

func (nc *NodeClient) Close() error {
	if nc.Conn != nil {
		return nc.Conn.Close()
	}
	return nil
}

// InitNodeClients initializes gRPC clients from config (legacy, returns empty slice).
func InitNodeClients(cfg *config.BackendConfig) ([]*NodeClient, error) {
	// Nodes are now loaded exclusively from database
	return []*NodeClient{}, nil
}

// DBNode represents a node loaded from database.
type DBNode struct {
	UUID                    string
	Name                    string
	Address                 string
	Port                    int
	APISchema               string
	APIPath                 string
	APIMetadata             string
	IsDisabled              bool
	ConsumptionMultiplier   int64
	IsTrafficTrackingActive bool
	TrafficResetDay         int
	TrafficLimitBytes       int64
	NotifyPercent           int
	ViewPosition            int
	CountryCode             string
	Tags                    []string
}

// LoadNodesFromDB loads all active nodes from the database.
func LoadNodesFromDB(manager *manager.DatabaseManager, cfg *config.BackendConfig) ([]DBNode, error) {
	var nodes []DBNode

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		query := `
			SELECT uuid, name, address, port, api_schema, api_path, api_metadata,
			       is_disabled, consumption_multiplier, is_traffic_tracking_active,
			       traffic_reset_day, traffic_limit_bytes, notify_percent,
			       view_position, country_code, tags
			FROM nodes
			WHERE is_disabled = 0
			ORDER BY view_position ASC, name ASC`

		rows, err := db.Query(query)
		if err != nil {
			return fmt.Errorf("query nodes: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var n DBNode
			var port sql.NullInt64
			var apiSchema, apiPath, apiMetadata, countryCode, tagsJSON sql.NullString
			var consumptionMultiplier, trafficLimitBytes, trafficResetDay, notifyPercent, viewPosition sql.NullInt64
			var isTrafficTrackingActive sql.NullBool

			err := rows.Scan(
				&n.UUID, &n.Name, &n.Address, &port, &apiSchema, &apiPath, &apiMetadata,
				&n.IsDisabled, &consumptionMultiplier, &isTrafficTrackingActive,
				&trafficResetDay, &trafficLimitBytes, &notifyPercent, &viewPosition,
				&countryCode, &tagsJSON,
			)
			if err != nil {
				return fmt.Errorf("scan node: %w", err)
			}

			// Handle nullable fields
			if port.Valid {
				n.Port = int(port.Int64)
			} else {
				n.Port = 8080
			}
			if apiSchema.Valid {
				n.APISchema = apiSchema.String
			} else {
				n.APISchema = "http"
			}
			if apiPath.Valid {
				n.APIPath = apiPath.String
			} else {
				n.APIPath = "/api"
			}
			if apiMetadata.Valid {
				n.APIMetadata = apiMetadata.String
			} else {
				n.APIMetadata = "{}"
			}
			if consumptionMultiplier.Valid {
				n.ConsumptionMultiplier = consumptionMultiplier.Int64
			} else {
				n.ConsumptionMultiplier = 100
			}
			if isTrafficTrackingActive.Valid {
				n.IsTrafficTrackingActive = isTrafficTrackingActive.Bool
			} else {
				n.IsTrafficTrackingActive = true
			}
			if trafficResetDay.Valid {
				n.TrafficResetDay = int(trafficResetDay.Int64)
			} else {
				n.TrafficResetDay = 1
			}
			if trafficLimitBytes.Valid {
				n.TrafficLimitBytes = trafficLimitBytes.Int64
			}
			if notifyPercent.Valid {
				n.NotifyPercent = int(notifyPercent.Int64)
			} else {
				n.NotifyPercent = 80
			}
			if viewPosition.Valid {
				n.ViewPosition = int(viewPosition.Int64)
			}
			if countryCode.Valid {
				n.CountryCode = countryCode.String
			}

			// Parse tags JSON
			if tagsJSON.Valid && tagsJSON.String != "" {
				if err := json.Unmarshal([]byte(tagsJSON.String), &n.Tags); err != nil {
					n.Tags = []string{}
				}
			} else {
				n.Tags = []string{}
			}

			nodes = append(nodes, n)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	cfg.Logger.Info("Loaded nodes from database", "count", len(nodes))
	return nodes, nil
}
