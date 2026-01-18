package config

type DisabledUsersConfigXray struct {
	Inbounds []XrayInbound `json:"inbounds"`
}

type XrayClient struct {
	Email    string `json:"email"`
	ID       string `json:"id,omitempty"`
	Password string `json:"password,omitempty"`
}

type XrayAccount struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type XraySettings struct {
	Clients    []XrayClient  `json:"clients,omitempty"`
	Accounts   []XrayAccount `json:"accounts,omitempty"`
	Auth       string        `json:"auth,omitempty"`
	Udp        bool          `json:"udp,omitempty"`
	Decryption *string       `json:"decryption,omitempty"`
	Address    *string       `json:"address,omitempty"`
}

type XrayInbound struct {
	Tag            string         `json:"tag"`
	Settings       XraySettings   `json:"settings"`
	Port           int            `json:"port"`
	Protocol       string         `json:"protocol"`
	Listen         string         `json:"listen"`
	StreamSettings map[string]any `json:"streamSettings,omitempty"`
	Sniffing       map[string]any `json:"sniffing,omitempty"`
	Allocate       map[string]any `json:"allocate,omitempty"`
}

type ConfigXray struct {
	Remarks          string           `json:"remarks,omitempty"`
	Log              map[string]any   `json:"log"`
	Dns              map[string]any   `json:"dns"`
	Routing          map[string]any   `json:"routing"`
	Inbounds         []XrayInbound    `json:"inbounds"`
	Outbounds        []map[string]any `json:"outbounds"`
	Policy           map[string]any   `json:"policy"`
	API              map[string]any   `json:"api"`
	Stats            map[string]any   `json:"stats"`
	FakeDNS          []map[string]any `json:"fakedns"`
	Reverse          map[string]any   `json:"reverse"`
	Transport        map[string]any   `json:"transport,omitempty"`
	Observatory      map[string]any   `json:"observatory,omitempty"`
	BurstObservatory map[string]any   `json:"burstObservatory,omitempty"`
}
