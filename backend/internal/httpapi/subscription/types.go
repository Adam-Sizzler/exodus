package subscription

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"exodus/internal/config"
	"exodus/internal/httpapi/subscriptionresponserules"
	"exodus/internal/httpapi/subscriptionsettings"
)

const (
	responseTypeBrowser    = "BROWSER"
	responseTypeXrayBase64 = "XRAY_BASE64"
	responseTypeXrayJSON   = "XRAY_JSON"
	responseTypeMihomo     = "MIHOMO"
	responseTypeStash      = "STASH"
	responseTypeClash      = "CLASH"
	responseTypeSingbox    = "SINGBOX"
	responseTypeBlock      = "BLOCK"
	responseTypeStatus404  = "STATUS_CODE_404"
	responseTypeStatus451  = "STATUS_CODE_451"
	responseTypeSocketDrop = "SOCKET_DROP"
)

var defaultResponseType = responseTypeXrayBase64

const defaultSubpageConfigUUID = "00000000-0000-0000-0000-000000000000"

type ExternalSquadOverrides struct {
	SubscriptionSettings *subscriptionsettings.SubscriptionSettings `json:"subscription_settings"`
	HostOverrides        map[string]HostOverride                    `json:"host_overrides"`
	ResponseHeaders      map[string]string                          `json:"response_headers"`
	HwidSettings         *HwidSettings                              `json:"hwid_settings"`
	CustomRemarks        *CustomRemarks                             `json:"custom_remarks"`
	Templates            map[string]string
}

type HwidSettingsInput struct {
	Enabled             *bool `json:"enabled"`
	FallbackDeviceLimit *int  `json:"fallbackDeviceLimit"`
	MaxDevicesAnnounce  *int  `json:"maxDevicesAnnounce"`
}

type UpdateExternalSquadInput struct {
	UUID              string           `json:"uuid"`
	Name              *string          `json:"name"`
	SubpageConfigUUID *string          `json:"subpageConfigUuid"`
	CustomRemarks     *json.RawMessage `json:"customRemarks"`
	HwidSettings      json.RawMessage  `json:"hwidSettings"`
}

type HostOverride struct {
	Address *string `json:"address"`
	Port    *int    `json:"port"`
	Remark  *string `json:"remark"`
	SNI     *string `json:"sni"`
	Host    *string `json:"host"`
	Path    *string `json:"path"`
}

type SubscriptionSettingsParsed struct {
	Raw subscriptionsettings.SubscriptionSettings

	HostOverrides         map[string]HostOverride
	ResponseHeaders       map[string]string
	CustomResponseHeaders map[string]string
	ResponseRules         *subscriptionresponserules.Config
	HwidSettings          HwidSettings
	CustomRemarks         CustomRemarks
}

type HwidSettings struct {
	Enabled             bool    `json:"enabled"`
	FallbackDeviceLimit int     `json:"fallbackDeviceLimit"`
	MaxDevicesAnnounce  *string `json:"maxDevicesAnnounce"`
}

type CustomRemarks struct {
	ExpiredUsers           []string `json:"expiredUsers"`
	LimitedUsers           []string `json:"limitedUsers"`
	DisabledUsers          []string `json:"disabledUsers"`
	EmptyHosts             []string `json:"emptyHosts"`
	HWIDMaxDevicesExceeded []string `json:"HWIDMaxDevicesExceeded"`
	HWIDNotSupported       []string `json:"HWIDNotSupported"`
}

type SubscriptionUser struct {
	TID                  int64
	UUID                 string
	ShortUUID            string
	Username             string
	Status               string
	TrafficLimitBytes    int64
	TrafficLimitStrategy string
	ExpireAt             time.Time
	TrojanPassword       string
	VlessUUID            string
	SSPassword           string
	NaivePassword        string
	ShadowtlsPassword    string
	Hysteria2Password    string
	AnytlsPassword       string
	HwidDeviceLimit      *int
	ExternalSquadUUID    *string
	UsedTrafficBytes     int64
	LifetimeUsedBytes    int64
}

type SubscriptionHost struct {
	UUID                         string
	ViewPosition                 int
	Remark                       string
	Address                      string
	Port                         int
	Path                         *string
	SNI                          *string
	Host                         *string
	ALPN                         *string
	Fingerprint                  *string
	SecurityLayer                string
	XHTTPExtraParams             *string
	MuxParams                    *string
	SingboxMuxParams             *string
	ClashMuxParams               *string
	SingboxCustomParams          *json.RawMessage
	MihomoCustomParams           *string
	SockoptParams                *string
	IsDisabled                   bool
	ServerDescription            *string
	OverrideProtocolCredential   bool
	ProtocolCredential           *string
	ShuffleHost                  bool
	MihomoX25519                 bool
	MihomoIPVersion              *string
	XrayJSONTemplateUUID         *string
	KeepSNIBlank                 bool
	ExcludeFromSubscriptionTypes []string
	Tag                          *string
	IsHidden                     bool
	OverrideSNIFromAddress       bool
	ConfigProfileUUID            *string
	ConfigProfileInboundUUID     *string

	InboundTag      *string
	InboundType     *string
	InboundNetwork  *string
	InboundSecurity *string
	InboundPort     *int
	InboundRaw      json.RawMessage
}

type RawHost struct {
	UUID        string  `json:"uuid"`
	Remark      string  `json:"remark"`
	Address     string  `json:"address"`
	Port        int     `json:"port"`
	Protocol    string  `json:"protocol"`
	Network     *string `json:"network,omitempty"`
	Security    *string `json:"security,omitempty"`
	Path        *string `json:"path,omitempty"`
	SNI         *string `json:"sni,omitempty"`
	Host        *string `json:"host,omitempty"`
	ALPN        *string `json:"alpn,omitempty"`
	Fingerprint *string `json:"fingerprint,omitempty"`
	IsDisabled  bool    `json:"is_disabled"`
	IsHidden    bool    `json:"is_hidden"`
}

type SubscriptionInfoUser struct {
	ShortUUID                string    `json:"shortUuid"`
	DaysLeft                 int       `json:"daysLeft"`
	TrafficUsed              string    `json:"trafficUsed"`
	TrafficLimit             string    `json:"trafficLimit"`
	LifetimeTrafficUsed      string    `json:"lifetimeTrafficUsed"`
	TrafficUsedBytes         string    `json:"trafficUsedBytes"`
	TrafficLimitBytes        string    `json:"trafficLimitBytes"`
	LifetimeTrafficUsedBytes string    `json:"lifetimeTrafficUsedBytes"`
	Username                 string    `json:"username"`
	ExpiresAt                time.Time `json:"expiresAt"`
	IsActive                 bool      `json:"isActive"`
	UserStatus               string    `json:"userStatus"`
	TrafficLimitStrategy     string    `json:"trafficLimitStrategy"`
}

type SubscriptionInfoResponse struct {
	IsFound         bool                 `json:"isFound"`
	User            SubscriptionInfoUser `json:"user"`
	Links           []string             `json:"links"`
	SSConfLinks     map[string]string    `json:"ssConfLinks"`
	SubscriptionURL string               `json:"subscriptionUrl"`
}

type RawSubscriptionResponse struct {
	User              map[string]interface{} `json:"user"`
	ConvertedUserInfo map[string]interface{} `json:"convertedUserInfo"`
	RawHosts          []RawHost              `json:"rawHosts"`
	Headers           map[string]string      `json:"headers"`
}

type SubscriptionHeaders struct {
	ContentDisposition     string            `json:"content-disposition,omitempty"`
	ProfileTitle           string            `json:"profile-title,omitempty"`
	ProfileUpdateInterval  string            `json:"profile-update-interval,omitempty"`
	ProfileWebPageURL      string            `json:"profile-web-page-url,omitempty"`
	SubscriptionUserInfo   string            `json:"subscription-userinfo,omitempty"`
	SubscriptionRefillDate string            `json:"subscription-refill-date,omitempty"`
	SupportURL             string            `json:"support-url,omitempty"`
	Announce               string            `json:"announce,omitempty"`
	Routing                string            `json:"routing,omitempty"`
	Extra                  map[string]string `json:"-"`
}

type SubscriptionWithConfig struct {
	Headers     map[string]string
	ContentType string
	Body        string
}

type HwidHeaders struct {
	Hwid        string
	Platform    *string
	OsVersion   *string
	DeviceModel *string
	UserAgent   *string
	RequestIP   *string
	Synthetic   bool
}

var (
	ErrHwidLimitExceeded = errors.New("hwid limit exceeded")
	ErrHwidRequired      = errors.New("hwid required")
	ErrUserDisabled      = errors.New("user disabled")
	ErrNoHosts           = errors.New("no hosts")
)

type XrayGenerator struct {
	cfg *config.BackendConfig
}

func NewXrayGenerator(cfg *config.BackendConfig) *XrayGenerator {
	return &XrayGenerator{cfg: cfg}
}

type MihomoGenerator struct {
	cfg *config.BackendConfig
}

func NewMihomoGenerator(cfg *config.BackendConfig) *MihomoGenerator {
	return &MihomoGenerator{cfg: cfg}
}

type SingboxGenerator struct {
	cfg *config.BackendConfig
}

func NewSingboxGenerator(cfg *config.BackendConfig) *SingboxGenerator {
	return &SingboxGenerator{cfg: cfg}
}

func (g *XrayGenerator) GenerateLinks(user SubscriptionUser, hosts []SubscriptionHost, settings SubscriptionSettingsParsed) ([]string, error) {
	links, _ := buildSubscriptionLinks(hosts, user)
	return links, nil
}

func (g *XrayGenerator) GenerateJSON(templateJSON []byte, user SubscriptionUser, hosts []SubscriptionHost, settings SubscriptionSettingsParsed) (string, error) {
	return generateXrayJSONConfig(templateJSON, hosts, user)
}

func (g *MihomoGenerator) Generate(templateYAML []byte, user SubscriptionUser, hosts []SubscriptionHost, settings SubscriptionSettingsParsed) (string, error) {
	return generateYAMLConfig(templateYAML, hosts, user)
}

func (g *SingboxGenerator) Generate(templateJSON []byte, user SubscriptionUser, hosts []SubscriptionHost, settings SubscriptionSettingsParsed) (string, error) {
	return generateSingboxConfig(templateJSON, hosts, user)
}

func matchResponseRules(rules *subscriptionresponserules.Config, headers http.Header) string {
	if rules == nil {
		return defaultResponseType
	}
	res := subscriptionresponserules.MatchRulesDetailed(rules, headers, "", func(s string) string { return s }, defaultResponseType)
	if res.Matched {
		return res.ResponseType
	}
	return defaultResponseType
}

func detectClientType(userAgent string) string {
	return inferPlatformFromUserAgent(userAgent)
}

