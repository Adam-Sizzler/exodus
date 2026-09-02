## Snippets

A snippet is a named piece of Sing-box configuration that is stored once and injected into any number of configuration profiles by name.

When you manage multiple profiles, every small change turns into a chore: for example, you have several profiles, and all of them share the same routing rules (blocking ads, torrent traffic, bypass lists). Without snippets, you have to repeat the edits in every profile. With snippets, you change the rule in one place, and it is picked up automatically by every profile that references it.

---

### Snippet format

A snippet is **always a JSON array of objects**, even if there is only one object inside.

```json
[
  {
    "type": "direct",
    "tag": "direct"
  }
]
```

Content requirements:

- it must be an array (`[]`), not an object (`{}`);
- the array must contain at least one item;
- empty objects (`{}`) inside the array are not allowed.

---

### Where snippets can be referenced

In Sing-box, a snippet can be referenced in four places. In three of them the target is an array, and the reference looks like an object `{ "snippet": "snippet-name" }`. In the fourth one the target is the configuration root, and the reference is an array of names in the `snippets` field.

#### 1. outbounds[]

Used to insert reusable outbound connection blocks (e.g. direct, block, socks, warp):

```json
{
  "outbounds": [
    {
      "snippet": "my-outbounds"
    }
  ]
}
```

#### 2. route.rules[]

Used to reuse Sing-box routing rules (blocking ads, bittorrent, geosite / geoip lists):

```json
{
  "route": {
    "rules": [
      {
        "snippet": "block-ads-and-torrents"
      }
    ]
  }
}
```

#### 3. endpoints[]

Used to insert Sing-box endpoint configurations (WireGuard, Tailscale, etc.):

```json
{
  "endpoints": [
    {
      "snippet": "my-endpoints"
    }
  ]
}
```

#### 4. Configuration root (snippets field)

At the configuration root, snippets are referenced by an array of names in the `snippets` field:

```json
{
  "snippets": ["dns-cloudflare", "log-preset"],
  "inbounds": [],
  "outbounds": []
}
```

Such a snippet must contain **root-level Sing-box sections** (e.g. `dns`, `log`, `experimental`). It is still an array of objects, where the keys are the names of Sing-box root sections:

```json
[
  {
    "dns": {
      "servers": [
        {
          "type": "udp",
          "tag": "dns-cloudflare",
          "server": "1.1.1.1"
        }
      ]
    }
  }
]
```

Sections can either be spread across separate array items or collected into a single object:

```json
[
  {
    "log": {
      "level": "info",
      "timestamp": true
    },
    "dns": {
      "servers": [
        {
          "type": "local",
          "tag": "dns-local"
        }
      ]
    }
  }
]
```

---

### How substitution works

#### In arrays

The object holding the reference `{ "snippet": "name" }` is replaced by the entire content of the snippet. If the snippet has several items, all of them take the place of that single reference.

**In the profile:**

```json
{
  "outbounds": [
    { "snippet": "default-fallback-outbounds" }
  ]
}
```

**The default-fallback-outbounds snippet:**

```json
[
  { "type": "block", "tag": "blocked" },
  { "type": "socks", "tag": "warp-fallback", "server": "127.0.0.1", "server_port": 40000 }
]
```

**Final config deployed to the node:**

```json
{
  "outbounds": [
    { "type": "block", "tag": "blocked" },
    { "type": "socks", "tag": "warp-fallback", "server": "127.0.0.1", "server_port": 40000 }
  ]
}
```

**The object holding the reference is replaced entirely.** If you write other keys next to `snippet`, they will be dropped.

If a snippet with the given name does not exist in the panel, the item holding the reference is removed from the array so it does not fail the Sing-box core launch.

#### At the root

All snippets listed in `snippets` are combined into a single set of sections and merged into the configuration root. The `snippets` field itself is removed from the final configuration.

Key rules:

1. **Sections already written in the profile are not overwritten.** If the profile has its own `log`, and a snippet also brings `log` — the one written in the profile stays. Hand-crafted configuration always takes precedence.
2. **System sections are never injected.** These are `inbounds`, `api`, `stats`, `metrics`, and `snippets` — they are managed by Exodus. If a snippet contains such sections, they are silently skipped.
3. **Order matters:** if two root snippets provide the same section, the one listed later in the array wins.

#### Processing order

The root is processed first, then the array elements. Because of that, a root snippet can bring in an `outbounds` section that contains references to other snippets — those will be expanded as well. Nesting is supported one level deep.
