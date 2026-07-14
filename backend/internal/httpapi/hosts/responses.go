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
		SecurityLayer:              rec.SecurityLayer,
		XHTTPExtraParams:           parseJSONAny(rec.XHTTPExtraParams),
		MuxParams:                  parseJSONAny(rec.MuxParams),
		SingboxMuxParams:           parseJSONAny(rec.SingboxMuxParams),
		ClashMuxParams:             rec.ClashMuxParams,
		SingboxCustomParams:        parseJSONAny(rec.SingboxCustomParams),
		MihomoCustomParams:         rec.MihomoCustomParams,
		SockoptParams:              parseJSONAny(rec.SockoptParams),
		FinalMask:                  parseJSONAny(rec.FinalMask),
		Inbound:                    HostInbound{ConfigProfileUUID: rec.ConfigProfileUUID, ConfigProfileInboundUUID: rec.ConfigProfileInboundUUID},
		ServerDescription:          rec.ServerDescription,
		Tags:                       ensureStringSlice(rec.Tags),
		IsHidden:                   rec.IsHidden,
		OverrideSNIFromAddress:     rec.OverrideSNIFromAddress,
		KeepSNIBlank:               rec.KeepSNIBlank,
		OverrideProtocolCredential: rec.OverrideProtocolCredential,
		ProtocolCredential:         rec.ProtocolCredential,
		VlessRouteID:               rec.VlessRouteID,
		PinnedPeerCertSha256:       rec.PinnedPeerCertSha256,
		VerifyPeerCertByName:       rec.VerifyPeerCertByName,
		AllowInsecure:              rec.AllowInsecure,
		ShuffleHost:                rec.ShuffleHost,
		MihomoX25519:               rec.MihomoX25519,
		MihomoIPVersion:            rec.MihomoIPVersion,
		Nodes:                      ensureStringSlice(nodes),
		XrayJSONTemplateUUID:       rec.XrayJSONTemplateUUID,
		ExcludedInternalSquads:     ensureStringSlice(excluded),
		ExcludeFromSubscription: func() []string {
			if len(rec.ExcludeTypes) == 0 {
				return []string{}
			}
			return rec.ExcludeTypes
		}(),
	}
}
