package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"v2ray-stat/logger"
	"v2ray-stat/backend/node/config"
	"v2ray-stat/proto"
)

func makeUserKey(username string) string {
	return username
}

func processXrayConfig(data []byte, userMap map[string]*proto.User, isDisabled bool, logger *logger.Logger) error {
	var cfg config.ConfigXray
	if err := json.Unmarshal(data, &cfg); err != nil {
		logger.Error("Failed to parse Xray config", "error", err)
		return fmt.Errorf("failed to parse config: %w", err)
	}

	for _, inbound := range cfg.Inbounds {
		switch inbound.Protocol {
		case "vless", "trojan":
			for _, client := range inbound.Settings.Clients {
				email := client.Email
				var id string
				if inbound.Protocol == "vless" {
					id = client.ID
				} else {
					id = client.Password
				}

				if id == "" || email == "" {
					logger.Debug("Skipping client without ID/Password or Email", "protocol", inbound.Protocol, "tag", inbound.Tag)
					continue
				}

				key := makeUserKey(email)
				if _, exists := userMap[key]; !exists {
					userMap[key] = &proto.User{
						Username:   email,
						IdInbounds: []*proto.IdInbound{},
						Enabled:    !isDisabled,
					}
				}
				userMap[key].IdInbounds = append(userMap[key].IdInbounds, &proto.IdInbound{
					Id:         id,
					InboundTag: inbound.Tag,
				})
				if isDisabled {
					userMap[key].Enabled = false
				}
			}

		case "mixed", "socks", "http":
			if inbound.Settings.Auth != "password" {
				logger.Debug("Skipping inbound without password auth", "protocol", inbound.Protocol, "tag", inbound.Tag)
				continue
			}
			for _, acc := range inbound.Settings.Accounts {
				user := acc.User
				pass := acc.Pass
				if user == "" || pass == "" {
					logger.Debug("Skipping account without User or Pass", "protocol", inbound.Protocol, "tag", inbound.Tag)
					continue
				}
				key := makeUserKey(user)
				if _, exists := userMap[key]; !exists {
					userMap[key] = &proto.User{
						Username:   user,
						IdInbounds: []*proto.IdInbound{},
						Enabled:    !isDisabled,
					}
				}
				userMap[key].IdInbounds = append(userMap[key].IdInbounds, &proto.IdInbound{
					Id:         pass,
					InboundTag: inbound.Tag,
				})
				if isDisabled {
					userMap[key].Enabled = false
				}
			}
		default:
			logger.Debug("Skipping inbound with unsupported protocol", "protocol", inbound.Protocol, "tag", inbound.Tag)
		}
	}
	return nil
}

func processSingboxConfig(data []byte, userMap map[string]*proto.User, isDisabled bool, logger *logger.Logger) error {
	var cfg config.ConfigSingbox
	if err := json.Unmarshal(data, &cfg); err != nil {
		logger.Error("Failed to parse Singbox config", "error", err)
		return fmt.Errorf("failed to parse config: %w", err)
	}

	for _, inbound := range cfg.Inbounds {
		switch inbound.Type {
		case "vless", "trojan", "mixed":
			for _, user := range inbound.Users {
				var id string
				switch inbound.Type {
				case "vless":
					id = user.UUID
				case "trojan", "mixed":
					id = user.Password
				}

				if id == "" || user.Name == "" {
					logger.Debug("Skipping user without ID/Password or Name", "type", inbound.Type, "tag", inbound.Tag)
					continue
				}

				key := makeUserKey(user.Name)
				if _, exists := userMap[key]; !exists {
					userMap[key] = &proto.User{
						Username:   user.Name,
						IdInbounds: []*proto.IdInbound{},
						Enabled:    !isDisabled,
					}
				}

				userMap[key].IdInbounds = append(userMap[key].IdInbounds, &proto.IdInbound{
					Id:         id,
					InboundTag: inbound.Tag,
				})
				if isDisabled {
					userMap[key].Enabled = false
				}
			}
		default:
			logger.Debug("Skipping inbound with unsupported type", "type", inbound.Type, "tag", inbound.Tag)
		}
	}
	return nil
}

// ListUsers retrieves a list of users from the node configuration.
func (s *NodeServer) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	userMap := make(map[string]*proto.User)
	mainConfigPath := s.Cfg.Core.Config
	disabledUsersPath := filepath.Join(s.Cfg.Core.Dir, ".disabled_users")

	s.Cfg.Logger.Debug("Starting ListUsers", "main_config_path", mainConfigPath, "disabled_users_path", disabledUsersPath)

	switch s.Cfg.Core.Type {
	case "xray":
		// Process main config
		data, err := os.ReadFile(mainConfigPath)
		if err != nil {
			s.Cfg.Logger.Error("Failed to read Xray main config", "path", mainConfigPath, "error", err)
			return nil, status.Errorf(codes.Internal, "read main config: %v", err)
		}
		err = processXrayConfig(data, userMap, false, s.Cfg.Logger)
		if err != nil {
			s.Cfg.Logger.Error("Failed to process Xray main config", "error", err)
			return nil, status.Errorf(codes.Internal, "process Xray main config: %v", err)
		}
		// Process disabled users
		disabledData, err := os.ReadFile(disabledUsersPath)
		if err == nil && len(disabledData) > 0 {
			err = processXrayConfig(disabledData, userMap, true, s.Cfg.Logger)
			if err != nil {
				s.Cfg.Logger.Error("Failed to process Xray disabled users", "path", disabledUsersPath, "error", err)
				return nil, status.Errorf(codes.Internal, "process Xray disabled users: %v", err)
			}
		} else if err != nil && !os.IsNotExist(err) {
			s.Cfg.Logger.Error("Failed to read Xray disabled users", "path", disabledUsersPath, "error", err)
			return nil, status.Errorf(codes.Internal, "read disabled users: %v", err)
		}

	case "singbox":
		// Process main config
		data, err := os.ReadFile(mainConfigPath)
		if err != nil {
			s.Cfg.Logger.Error("Failed to read Singbox main config", "path", mainConfigPath, "error", err)
			return nil, status.Errorf(codes.Internal, "read main config: %v", err)
		}
		err = processSingboxConfig(data, userMap, false, s.Cfg.Logger)
		if err != nil {
			s.Cfg.Logger.Error("Failed to process Singbox main config", "error", err)
			return nil, status.Errorf(codes.Internal, "process Singbox main config: %v", err)
		}
		// Process disabled users
		disabledData, err := os.ReadFile(disabledUsersPath)
		if err == nil && len(disabledData) > 0 {
			err = processSingboxConfig(disabledData, userMap, true, s.Cfg.Logger)
			if err != nil {
				s.Cfg.Logger.Error("Failed to process Singbox disabled users", "path", disabledUsersPath, "error", err)
				return nil, status.Errorf(codes.Internal, "process Singbox disabled users: %v", err)
			}
		} else if err != nil && !os.IsNotExist(err) {
			s.Cfg.Logger.Error("Failed to read Singbox disabled users", "path", disabledUsersPath, "error", err)
			return nil, status.Errorf(codes.Internal, "read disabled users: %v", err)
		}
	default:
		s.Cfg.Logger.Error("Unsupported core type", "type", s.Cfg.Core.Type)
		return nil, status.Errorf(codes.InvalidArgument, "unsupported core type: %s", s.Cfg.Core.Type)
	}

	resp := &proto.ListUsersResponse{
		Users: make([]*proto.User, 0, len(userMap)),
	}
	for _, user := range userMap {
		resp.Users = append(resp.Users, user)
	}
	s.Cfg.Logger.Debug("Returning users", "count", len(resp.Users))
	return resp, nil
}
