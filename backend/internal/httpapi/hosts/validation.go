package hosts

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"exodus/internal/util"

	"github.com/google/uuid"
)

func validateCreateRequest(req HostCreateRequestAPI) error {
	if utf8.RuneCountInString(req.Remark) < 1 {
		return fmt.Errorf("remark must be at least 1 character")
	}
	if utf8.RuneCountInString(req.Remark) > 100 {
		return fmt.Errorf("remark must be less than 100 characters")
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
	if err := validateHostTags(req.Tags); err != nil {
		return err
	}
	if req.ServerDescription != nil && len(*req.ServerDescription) > 30 {
		return fmt.Errorf("serverDescription must be less than 30 characters")
	}
	if err := validateVlessRouteID(req.VlessRouteID); err != nil {
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
	if err := validateMihomoIPVersion(req.MihomoIPVersion); err != nil {
		return err
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

func validateUpdateRequest(req hostUpdateFields) error {
	if req.Remark.Set {
		if req.Remark.Value == nil {
			return fmt.Errorf("remark cannot be null")
		}
		if utf8.RuneCountInString(*req.Remark.Value) < 1 {
			return fmt.Errorf("remark must be at least 1 character")
		}
		if utf8.RuneCountInString(*req.Remark.Value) > 100 {
			return fmt.Errorf("remark must be less than 100 characters")
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
	if req.Tags != nil {
		if err := validateHostTags(req.Tags); err != nil {
			return err
		}
	}
	if req.ServerDescription.Set && req.ServerDescription.Value != nil {
		if len(*req.ServerDescription.Value) > 30 {
			return fmt.Errorf("serverDescription must be less than 30 characters")
		}
	}
	if req.VlessRouteID.Set {
		if err := validateVlessRouteID(req.VlessRouteID.Value); err != nil {
			return err
		}
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
	if req.MihomoIPVersion.Set {
		if err := validateMihomoIPVersion(req.MihomoIPVersion.Value); err != nil {
			return err
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

func validateHostTags(tags []string) error {
	if len(tags) > 10 {
		return fmt.Errorf("maximum 10 tags")
	}
	for _, tag := range tags {
		if tag == "" {
			return fmt.Errorf("tag cannot be empty")
		}
		if !hostTagRegex.MatchString(tag) {
			return fmt.Errorf("invalid tag format")
		}
		if len(tag) > 36 {
			return fmt.Errorf("tag must be less than 36 characters")
		}
	}
	return nil
}

func validateUUIDList(values []string) error {
	return util.ValidateUUIDsAllowEmpty(values)
}

func validateTemplateTypes(values []string) error {
	for _, value := range values {
		if _, ok := allowedTemplateTypes[value]; !ok {
			return fmt.Errorf("invalid subscription template type")
		}
	}
	return nil
}

func validateMihomoIPVersion(value *string) error {
	if value == nil {
		return nil
	}
	normalized := strings.ToLower(strings.TrimSpace(*value))
	if normalized == "" {
		return fmt.Errorf("invalid mihomoIpVersion")
	}
	if _, ok := allowedMihomoIPVersions[normalized]; !ok {
		return fmt.Errorf("invalid mihomoIpVersion")
	}
	return nil
}

func validateVlessRouteID(value *int) error {
	if value == nil {
		return nil
	}
	if *value < 0 || *value > 65535 {
		return fmt.Errorf("invalid vlessRouteId")
	}
	return nil
}
