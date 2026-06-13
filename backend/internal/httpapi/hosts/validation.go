package hosts

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func validateCreateRequest(req HostCreateRequestAPI) error {
	if len(req.Remark) < 1 {
		return fmt.Errorf("remark must be at least 1 character")
	}
	if len(req.Remark) > 40 {
		return fmt.Errorf("remark must be less than 40 characters")
	}
	if req.Port < 1 || req.Port > 65535 {
		return fmt.Errorf("invalid port")
	}
	if req.Inbound.ConfigProfileUUID == nil || req.Inbound.ConfigProfileInboundUUID == nil {
		return fmt.Errorf("inbound configProfileUuid and configProfileInboundUuid are required")
	}
	if _, err := uuid.Parse(*req.Inbound.ConfigProfileUUID); err != nil {
		return fmt.Errorf("invalid configProfileUuid")
	}
	if _, err := uuid.Parse(*req.Inbound.ConfigProfileInboundUUID); err != nil {
		return fmt.Errorf("invalid configProfileInboundUuid")
	}
	if req.Tag != nil && *req.Tag != "" {
		if !hostTagRegex.MatchString(*req.Tag) {
			return fmt.Errorf("invalid tag format")
		}
		if len(*req.Tag) > 32 {
			return fmt.Errorf("tag must be less than 32 characters")
		}
	} else if req.Tag != nil && *req.Tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if req.ServerDescription != nil && len(*req.ServerDescription) > 30 {
		return fmt.Errorf("serverDescription must be less than 30 characters")
	}
	if err := validateProtocolCredentialCreate(req.OverrideProtocolCredential, req.ProtocolCredential); err != nil {
		return err
	}
	if req.SecurityLayer != nil {
		if strings.TrimSpace(*req.SecurityLayer) == "" {
			return fmt.Errorf("invalid securityLayer")
		}
		if _, ok := allowedSecurityLayers[strings.ToUpper(*req.SecurityLayer)]; !ok {
			return fmt.Errorf("invalid securityLayer")
		}
	}
	if req.ALPN != nil && *req.ALPN != "" {
		if _, ok := allowedAlpn[*req.ALPN]; !ok {
			return fmt.Errorf("invalid alpn")
		}
	} else if req.ALPN != nil && *req.ALPN == "" {
		return fmt.Errorf("invalid alpn")
	}
	if req.Fingerprint != nil && *req.Fingerprint != "" {
		if _, ok := allowedFingerprints[*req.Fingerprint]; !ok {
			return fmt.Errorf("invalid fingerprint")
		}
	} else if req.Fingerprint != nil && *req.Fingerprint == "" {
		return fmt.Errorf("invalid fingerprint")
	}
	if req.XrayJSONTemplateUUID != nil && *req.XrayJSONTemplateUUID != "" {
		if _, err := uuid.Parse(*req.XrayJSONTemplateUUID); err != nil {
			return fmt.Errorf("invalid xrayJsonTemplateUuid")
		}
	} else if req.XrayJSONTemplateUUID != nil && *req.XrayJSONTemplateUUID == "" {
		return fmt.Errorf("invalid xrayJsonTemplateUuid")
	}
	if err := validateUUIDList(req.Nodes); err != nil {
		return err
	}
	if err := validateUUIDList(req.ExcludedInternalSquads); err != nil {
		return err
	}
	if err := validateTemplateTypes(req.ExcludeFromSubscription); err != nil {
		return err
	}
	return nil
}

func validateUpdateRequest(req HostUpdateRequestAPI) error {
	if req.Remark.Set {
		if req.Remark.Value == nil {
			return fmt.Errorf("remark cannot be null")
		}
		if len(*req.Remark.Value) > 40 {
			return fmt.Errorf("remark must be less than 40 characters")
		}
	}
	if req.Address.Set && req.Address.Value == nil {
		return fmt.Errorf("address cannot be null")
	}
	if req.Path.Set && req.Path.Value == nil {
		return fmt.Errorf("path cannot be null")
	}
	if req.SNI.Set && req.SNI.Value == nil {
		return fmt.Errorf("sni cannot be null")
	}
	if req.Host.Set && req.Host.Value == nil {
		return fmt.Errorf("host cannot be null")
	}
	if req.Port != nil && (*req.Port < 1 || *req.Port > 65535) {
		return fmt.Errorf("invalid port")
	}
	if req.Tag.Set {
		if req.Tag.Value == nil {
			// nullable
		} else if *req.Tag.Value == "" {
			return fmt.Errorf("tag cannot be empty")
		} else if !hostTagRegex.MatchString(*req.Tag.Value) {
			return fmt.Errorf("invalid tag format")
		} else if len(*req.Tag.Value) > 32 {
			return fmt.Errorf("tag must be less than 32 characters")
		}
	}
	if req.ServerDescription.Set && req.ServerDescription.Value != nil {
		if len(*req.ServerDescription.Value) > 30 {
			return fmt.Errorf("serverDescription must be less than 30 characters")
		}
	}
	if err := validateProtocolCredentialUpdate(req.OverrideProtocolCredential, req.ProtocolCredential); err != nil {
		return err
	}
	if req.SecurityLayer != nil {
		if strings.TrimSpace(*req.SecurityLayer) == "" {
			return fmt.Errorf("invalid securityLayer")
		}
		if _, ok := allowedSecurityLayers[strings.ToUpper(*req.SecurityLayer)]; !ok {
			return fmt.Errorf("invalid securityLayer")
		}
	}
	if req.ALPN.Set {
		if req.ALPN.Value == nil {
			// nullable
		} else if *req.ALPN.Value == "" {
			return fmt.Errorf("invalid alpn")
		} else if _, ok := allowedAlpn[*req.ALPN.Value]; !ok {
			return fmt.Errorf("invalid alpn")
		}
	}
	if req.Fingerprint.Set {
		if req.Fingerprint.Value == nil {
			// nullable
		} else if *req.Fingerprint.Value == "" {
			return fmt.Errorf("invalid fingerprint")
		} else if _, ok := allowedFingerprints[*req.Fingerprint.Value]; !ok {
			return fmt.Errorf("invalid fingerprint")
		}
	}
	if req.XrayJSONTemplateUUID.Set {
		if req.XrayJSONTemplateUUID.Value == nil {
			// nullable
		} else if *req.XrayJSONTemplateUUID.Value == "" {
			return fmt.Errorf("invalid xrayJsonTemplateUuid")
		} else if _, err := uuid.Parse(*req.XrayJSONTemplateUUID.Value); err != nil {
			return fmt.Errorf("invalid xrayJsonTemplateUuid")
		}
	}
	if req.Inbound != nil {
		if req.Inbound.ConfigProfileUUID == nil || req.Inbound.ConfigProfileInboundUUID == nil {
			return fmt.Errorf("inbound configProfileUuid and configProfileInboundUuid are required")
		}
		if _, err := uuid.Parse(*req.Inbound.ConfigProfileUUID); err != nil {
			return fmt.Errorf("invalid configProfileUuid")
		}
		if _, err := uuid.Parse(*req.Inbound.ConfigProfileInboundUUID); err != nil {
			return fmt.Errorf("invalid configProfileInboundUuid")
		}
	}
	if err := validateUUIDList(req.Nodes); err != nil {
		return err
	}
	if err := validateUUIDList(req.ExcludedInternalSquads); err != nil {
		return err
	}
	if err := validateTemplateTypes(req.ExcludeFromSubscription); err != nil {
		return err
	}
	return nil
}

func validateUUIDList(values []string) error {
	for _, value := range values {
		if _, err := uuid.Parse(value); err != nil {
			return fmt.Errorf("invalid uuid value")
		}
	}
	return nil
}

func validateTemplateTypes(values []string) error {
	for _, value := range values {
		if _, ok := allowedTemplateTypes[value]; !ok {
			return fmt.Errorf("invalid subscription template type")
		}
	}
	return nil
}

func validateProtocolCredentialCreate(override *bool, value *string) error {
	if err := validateProtocolCredentialValue(value); err != nil {
		return err
	}
	if coalesceBool(override, false) && normalizeProtocolCredentialPointer(value) == nil {
		return fmt.Errorf("protocolCredential is required when overrideProtocolCredential is enabled")
	}
	return nil
}

func validateProtocolCredentialUpdate(override *bool, value OptionalString) error {
	if value.Set {
		if err := validateProtocolCredentialValue(value.Value); err != nil {
			return err
		}
	}
	if override != nil && *override && value.Set && normalizeProtocolCredentialPointer(value.Value) == nil {
		return fmt.Errorf("protocolCredential is required when overrideProtocolCredential is enabled")
	}
	return nil
}

func validateProtocolCredentialValue(value *string) error {
	if value == nil {
		return nil
	}
	if len(strings.TrimSpace(*value)) > maxProtocolCredentialLength {
		return fmt.Errorf("protocolCredential must be less than 256 characters")
	}
	return nil
}
