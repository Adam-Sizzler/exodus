# Node Plugins
 
Node plugins are optional modules that can be activated on Exodus Node: **Ingress Filter**, **Egress Filter**, **HAProxy Auth**, **Shared Lists**, and **Executor** action commands.
 
> **Important:** Exercise caution when configuring network filter plugins. Ingress Filter and Egress Filter interact directly with the host's `nftables` firewall.
 
---
 
## Requirements
 
1. `cap_add: NET_ADMIN` capability in the node `docker-compose.yml`.
2. Linux Kernel **5.7** or higher (`uname -r`).
3. `nftables` utility installed on host (`nft --version`).
 
```yaml title="docker-compose.yml"
services:
    exodus-node:
        container_name: exodus-node
        image: exodus/node:latest
        network_mode: host
        restart: always
        cap_add:
            - NET_ADMIN
        ulimits:
            nofile:
                soft: 1048576
                hard: 1048576
```
 
---
 
## nftables Rule Tables
 
When network filter plugins are enabled, Exodus Node automatically provisions the following `nftables` tables:
 
```nftables
table ip exodus {
        counter processed {
                packets 32 bytes 2060
        }
 
        counter ingress-filter-ip {
                packets 0 bytes 0
        }
 
        counter egress-filter-ip {
                packets 0 bytes 0
        }
 
        counter egress-filter-port {
                packets 0 bytes 0
        }
 
        set ingress-filter-ip {
                type ipv4_addr
                flags timeout
                counter
        }
 
        set egress-filter-ip {
                type ipv4_addr
                flags timeout
                counter
        }
 
        set egress-filter-port {
                type inet_proto . inet_service
                flags timeout
                counter
        }
 
        chain input {
                type filter hook input priority filter - 10; policy accept;
                counter name "processed"
                ip saddr @ingress-filter-ip log prefix "ingress-filter-ip: " counter name "ingress-filter-ip" drop
        }
 
        chain forward {
                type filter hook forward priority filter - 10; policy accept;
                counter name "processed"
                ip saddr @ingress-filter-ip log prefix "ingress-filter-ip: " counter name "ingress-filter-ip" drop
        }
 
        chain output {
                type filter hook output priority filter - 10; policy accept;
                ip daddr @egress-filter-ip counter name "egress-filter-ip" drop
                meta l4proto . th dport @egress-filter-port counter name "egress-filter-port" drop
        }
}
 
table ip6 exodus6 {
        counter processed {
                packets 0 bytes 0
        }
 
        counter ingress-filter-ip6 {
                packets 0 bytes 0
        }
 
        counter egress-filter-ip6 {
                packets 0 bytes 0
        }
 
        counter egress-filter-port6 {
                packets 0 bytes 0
        }
 
        set ingress-filter-ip6 {
                type ipv6_addr
                flags timeout
                counter
        }
 
        set egress-filter-ip6 {
                type ipv6_addr
                flags timeout
                counter
        }
 
        set egress-filter-port6 {
                type inet_proto . inet_service
                flags timeout
                counter
        }
 
        chain input {
                type filter hook input priority filter - 10; policy accept;
                counter name "processed"
                ip6 saddr @ingress-filter-ip6 log prefix "ingress-filter-ip: " counter name "ingress-filter-ip6" drop
        }
 
        chain forward {
                type filter hook forward priority filter - 10; policy accept;
                counter name "processed"
                ip6 saddr @ingress-filter-ip6 log prefix "ingress-filter-ip: " counter name "ingress-filter-ip6" drop
        }
 
        chain output {
                type filter hook output priority filter - 10; policy accept;
                ip6 daddr @egress-filter-ip6 counter name "egress-filter-ip6" drop
                meta l4proto . th dport @egress-filter-port6 counter name "egress-filter-port6" drop
        }
}
```
 
---
 
## CIDR Support
 
CIDR notation (IPv4 & IPv6) is supported in **Ingress Filter**, **Egress Filter**, and **Shared Lists**.
 
Examples: `192.168.1.1`, `10.0.0.0/8`, `172.16.0.0/12`, `2001:db8::/32`.
 
---
 
## JSON Configuration Schema
 
```json
{
    "ingressFilter": {
        "enabled": false,
        "blockedIps": []
    },
    "egressFilter": {
        "enabled": false,
        "blockedIps": [],
        "blockedPorts": []
    },
    "haproxyAuth": {
        "enabled": false,
        "inboundTags": ["*"]
    }
}
```
 
---
 
## Ingress Filter
 
Permanently drops incoming connections from specified IP addresses or CIDR subnets using `nftables`.
 
```json
"ingressFilter": {
    "enabled": true,
    "blockedIps": ["192.168.1.1", "10.0.0.0/8", "ext:blocked-ips"]
}
```
 
---
 
## Egress Filter
 
Blocks outbound connections from the node to target IP addresses or destination ports (e.g. SMTP ports 25, 465, 587).
 
```json
"egressFilter": {
    "enabled": true,
    "blockedIps": ["10.0.0.1", "ext:internal-networks"],
    "blockedPorts": [25, 465, 587]
}
```
 
---
 
## HAProxy Auth Module
 
Generates and maintains a CSV user credentials file (`/opt/app/haproxy/data/users.csv`) inside the node container for fast Lua-based authentication in HAProxy with zero-downtime socket reload.
 
```json
"haproxyAuth": {
    "enabled": true,
    "inboundTags": ["*"]
}
```
 
- Specify `["*"]` in `inboundTags` to include all supported inbounds on the node, or list specific inbound tags (e.g. `["vless-inbound", "trojan-inbound", "naive-inbound"]`).
- Automatically updates `/opt/app/haproxy/data/users.csv` on the node.
- Format: `username,credential` (without legacy 1/0 prefix).
- Supports: `VLESS UUID`, SHA-224 hash for `Trojan credentials`, SHA-256 hash for `AnyTLS credentials`, and `basic:<base64(user:pass)>` for `NaiveProxy`.
 
---
 
## Shared Lists
 
Shared lists are managed in the dedicated **Shared Lists** menu, creating globally reusable collections of IP addresses, CIDR ranges, or Autonomous System Numbers (ASN).
 
Reference them in plugins by prefixing with `ext:list_name`.
 
Example in Ingress Filter:
```json
"ingressFilter": {
    "enabled": true,
    "blockedIps": ["ext:my-blocked-subnets"]
}
```
 
---
 
## Executor Actions
 
The **Executor** menu in the header allows executing management commands:
- **Block IPs**: Temporarily block IP addresses in `nftables` with a custom duration.
- **Unblock IPs**: Remove temporary IP blocks.
- **Reset nftables**: Force recreation of `nftables` rule tables.
