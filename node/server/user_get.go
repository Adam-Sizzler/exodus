package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"v2ray-stat/node/config"
	"v2ray-stat/proto"
)

// makeUserKey creates a unique key for a user based on their username.
func makeUserKey(username string) string {
	return username
}

// processUsers processes users from configuration and adds them to userMap.
func processUsers[T any](data []byte, userMap map[string]*proto.User, cfg *config.NodeConfig, isDisabled bool, process func(cfg T, userMap map[string]*proto.User, isDisabled bool)) error {
	var config T
	if err := json.Unmarshal(data, &config); err != nil {
		cfg.Logger.Error("Failed to parse config", "error", err)
		return fmt.Errorf("failed to parse config: %w", err)
	}
	process(config, userMap, isDisabled)
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
		err = processUsers(data, userMap, s.Cfg, false, func(cfg config.ConfigXray, userMap map[string]*proto.User, isDisabled bool) {
			for _, inbound := range cfg.Inbounds {
				switch inbound.Protocol {
				case "vless":
					for _, client := range inbound.Settings.Clients {
						key := makeUserKey(client.Email)
						if _, exists := userMap[key]; !exists {
							userMap[key] = &proto.User{
								Username:   client.Email,
								IdInbounds: []*proto.IdInbound{},
								Enabled:    !isDisabled,
							}
						}
						userMap[key].IdInbounds = append(userMap[key].IdInbounds, &proto.IdInbound{
							Id:         client.Password,
							InboundTag: inbound.Tag,
						})
					}
				case "trojan":
					for _, client := range inbound.Settings.Clients {
						key := makeUserKey(client.Email)
						if _, exists := userMap[key]; !exists {
							userMap[key] = &proto.User{
								Username:   client.Email,
								IdInbounds: []*proto.IdInbound{},
								Enabled:    !isDisabled,
							}
						}
						userMap[key].IdInbounds = append(userMap[key].IdInbounds, &proto.IdInbound{
							Id:         client.Password,
							InboundTag: inbound.Tag,
						})
					}
				case "mixed", "socks", "http":
					if inbound.Settings.Auth != "password" {
						s.Cfg.Logger.Debug("Skipping inbound with no password auth", "protocol", inbound.Protocol, "tag", inbound.Tag)
						continue
					}
					for _, acc := range inbound.Settings.Accounts {
						key := makeUserKey(acc.User)
						if _, exists := userMap[key]; !exists {
							userMap[key] = &proto.User{
								Username:   acc.User,
								IdInbounds: []*proto.IdInbound{},
								Enabled:    !isDisabled,
							}
						}
						userMap[key].IdInbounds = append(userMap[key].IdInbounds, &proto.IdInbound{
							Id:         acc.Pass,
							InboundTag: inbound.Tag,
						})
					}
				default:
					s.Cfg.Logger.Debug("Skipping inbound with unsupported protocol", "protocol", inbound.Protocol, "tag", inbound.Tag)
					continue
				}
			}
		})
		if err != nil {
			s.Cfg.Logger.Error("Failed to process Xray main config", "error", err)
			return nil, status.Errorf(codes.Internal, "process Xray main config: %v", err)
		}

		// Process disabled users
		disabledData, err := os.ReadFile(disabledUsersPath)
		if err == nil && len(disabledData) > 0 {
			err = processUsers(disabledData, userMap, s.Cfg, true, func(cfg config.DisabledUsersConfigXray, userMap map[string]*proto.User, isDisabled bool) {
				for _, inbound := range cfg.Inbounds {
					switch inbound.Protocol {
					case "vless":
						for _, client := range inbound.Settings.Clients {
							key := makeUserKey(client.Email)
							if _, exists := userMap[key]; !exists {
								userMap[key] = &proto.User{
									Username:   client.Email,
									IdInbounds: []*proto.IdInbound{},
									Enabled:    !isDisabled,
								}
							}
							userMap[key].IdInbounds = append(userMap[key].IdInbounds, &proto.IdInbound{
								Id:         client.ID,
								InboundTag: inbound.Tag,
							})
							userMap[key].Enabled = !isDisabled
						}
					case "trojan":
						for _, client := range inbound.Settings.Clients {
							key := makeUserKey(client.Email)
							if _, exists := userMap[key]; !exists {
								userMap[key] = &proto.User{
									Username:   client.Email,
									IdInbounds: []*proto.IdInbound{},
									Enabled:    !isDisabled,
								}
							}
							userMap[key].IdInbounds = append(userMap[key].IdInbounds, &proto.IdInbound{
								Id:         client.Password,
								InboundTag: inbound.Tag,
							})
							userMap[key].Enabled = !isDisabled
						}
					case "mixed", "socks", "http":
						if inbound.Settings.Auth != "password" {
							s.Cfg.Logger.Debug("Skipping disabled inbound with no password auth", "protocol", inbound.Protocol, "tag", inbound.Tag)
							continue
						}
						for _, acc := range inbound.Settings.Accounts {
							key := makeUserKey(acc.User)
							if _, exists := userMap[key]; !exists {
								userMap[key] = &proto.User{
									Username:   acc.User,
									IdInbounds: []*proto.IdInbound{},
									Enabled:    !isDisabled,
								}
							}
							userMap[key].IdInbounds = append(userMap[key].IdInbounds, &proto.IdInbound{
								Id:         acc.Pass,
								InboundTag: inbound.Tag,
							})
							userMap[key].Enabled = !isDisabled
						}
					default:
						s.Cfg.Logger.Debug("Skipping disabled inbound with unsupported protocol", "protocol", inbound.Protocol, "tag", inbound.Tag)
						continue
					}
				}
			})
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
		err = processUsers(data, userMap, s.Cfg, false, func(cfg config.ConfigSingbox, userMap map[string]*proto.User, isDisabled bool) {
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
						if id == "" {
							s.Cfg.Logger.Debug("Skipping user without ID/Password", "type", inbound.Type, "tag", inbound.Tag, "name", user.Name)
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
					}
				default:
					s.Cfg.Logger.Debug("Skipping inbound with unsupported type", "type", inbound.Type, "tag", inbound.Tag)
					continue
				}
			}
		})
		if err != nil {
			s.Cfg.Logger.Error("Failed to process Singbox main config", "error", err)
			return nil, status.Errorf(codes.Internal, "process Singbox main config: %v", err)
		}

		// Process disabled users
		disabledData, err := os.ReadFile(disabledUsersPath)
		if err == nil && len(disabledData) > 0 {
			err = processUsers(disabledData, userMap, s.Cfg, true, func(cfg config.DisabledUsersConfigSingbox, userMap map[string]*proto.User, isDisabled bool) {
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
							if id == "" {
								s.Cfg.Logger.Debug("Skipping disabled user without ID/Password", "type", inbound.Type, "tag", inbound.Tag, "name", user.Name)
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
							userMap[key].Enabled = !isDisabled
						}
					default:
						s.Cfg.Logger.Debug("Skipping disabled inbound with unsupported type", "type", inbound.Type, "tag", inbound.Tag)
						continue
					}
				}
			})
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
