export const BASIC_MUX_PARAMS = `{
  "enabled": true,
  "concurrency": -1,
  "xudpConcurrency": 16,
  "xudpProxyUDP443": "reject"
}`

export const BASIC_SINGBOX_MUX_PARAMS = `{
  "enabled": true,
  "padding": true
}`

export const BASIC_CLASH_MUX_PARAMS = `enabled: true
protocol: smux
max-connections: 4
min-streams: 4`

export const BASIC_SOCKOPT_PARAMS = `{
  "mark": 0,
  "tcpMaxSeg": 1440,
  "tcpFastOpen": false,
  "tproxy": "off",
  "domainStrategy": "AsIs",
  "dialerProxy": "",
  "happyEyeballs": {},
  "acceptProxyProtocol": false,
  "tcpKeepAliveInterval": 0,
  "tcpKeepAliveIdle": 300,
  "tcpUserTimeout": 10000,
  "tcpcongestion": "bbr",
  "interface": "wg0",
  "V6Only": false,
  "tcpWindowClamp": 600,
  "tcpMptcp": false,
  "tcpNoDelay": false
}`
