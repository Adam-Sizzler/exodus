package hosts

func mapHostRecordToAPI(rec hostRecord, nodes []string, excluded []string) HostAPI {
	return HostAPI{
		UUID:                       rec.UUID,
		ViewPosition:               rec.ViewPosition,
		Remark:                     rec.Remark,
		Address:                    rec.Address,
		Port:                       rec.Port,
		Path:                       rec.Path,
		SNI:                        rec.SNI,
		Host:                       rec.Host,
		ALPN:                       rec.ALPN,
		Fingerprint:                rec.Fingerprint,
		IsDisabled:                 rec.IsDisabled,
		SecurityLayer:           rec.SecurityLayer,
		XHTTPExtraParams:        parseJSONAny(rec.XHTTPExtraParams),
		MuxParams:               parseJSONAny(rec.MuxParams),
		SockoptParams:           parseJSONAny(rec.SockoptParams),
		FinalMask:               parseJSONAny(rec.FinalMask),
		Inbound:                 HostInbound{ConfigProfileUUID: rec.ConfigProfileUUID, ConfigProfileInboundUUID: rec.ConfigProfileInboundUUID},
		ServerDescription:       rec.ServerDescription,
		Tags:                    ensureStringSlice(rec.Tags),
		IsHidden:                rec.IsHidden,
		OverrideSNIFromAddress:  rec.OverrideSNIFromAddress,
		KeepSNIBlank:            rec.KeepSNIBlank,
		VlessRouteID:            rec.VlessRouteID,
		PinnedPeerCertSha256:    rec.PinnedPeerCertSha256,
		VerifyPeerCertByName:    rec.VerifyPeerCertByName,
		ShuffleHost:             rec.ShuffleHost,
		MihomoX25519:            rec.MihomoX25519,
		MihomoIPVersion:         rec.MihomoIPVersion,
		Nodes:                   ensureStringSlice(nodes),
		XrayJSONTemplateUUID:    rec.XrayJSONTemplateUUID,
		ExcludedInternalSquads:  ensureStringSlice(excluded),
		ExcludeFromSubscription: func() []string {
			if len(rec.ExcludeTypes) == 0 {
				return []string{}
			}
			return rec.ExcludeTypes
		}(),
		Mapper: func() interface{} {
			m := parseJSONAny(rec.Mapper)
			if m == nil {
				return map[string]interface{}{}
			}
			return m
		}(),
	}
}
