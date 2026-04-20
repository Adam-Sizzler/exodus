## Module: HAPROXY

In short: this module builds and maintains the credentials table used by HAProxy Lua auth during node deploy.

### What this module does

- Creates and updates **one** file with users.
- Includes all users of the current node, without binding them to a specific inbound.
- Writes each user's `VLESS UUID`.
- Writes each user's `Trojan credential` as a SHA-224 hash, 56 characters long.
- Sets the enabled state to `1`.

### File path

- Relative path in module logic: `haproxy/data/users.csv`
- Absolute path inside the node container: `/app/haproxy/data/users.csv`

### Row format

`1,username,credential`

### Example

```text
1,pablo_escobar,90cdd126-819d-5f2f-a1dd-4cbf085182b9
1,alvarez,9f37a85cafb66ab9c36f866e5c787b2e42a997716c02337f532a259c
1,guest,8f49d237686bee189aaedaba1ad434c5580356d8e8b5a3d702bd0228
```

### Why it is needed

- `haproxy/lua/auth.lua` loads users from `/etc/haproxy/data/users.csv` and detects the protocol.
- Trojan traffic is detected by credential, including the MUX check.
- VLESS traffic is detected by UUID from the binary handshake.
- `haproxy/lua/users_loader.lua` parses the CSV into lookup maps.
- `trojan[sha224] -> username`
- `vless[uuid] -> username`
- Then `haproxy.cfg` uses `lua.identify_protocol` and routes traffic to the `trojan` or `vless` backend.

### Important

- The module does not modify `haproxy.cfg` or Lua scripts.
- The module manages only `users.csv`.
