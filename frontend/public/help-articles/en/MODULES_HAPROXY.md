## Module: HAPROXY
 
In short: this module builds and maintains the credentials table used by HAProxy Lua auth during node deploy with hot reload via runtime socket.
 
### What this module does
 
- Creates and updates **one** user credentials file `users.csv`.
- Includes active users of the current node.
- Writes each user's `VLESS UUID`.
- Writes each user's `Trojan credential` as a SHA-224 hash (56 characters hex).
- Writes each user's `AnyTLS credential` as a SHA-256 hash (64 characters hex).
- Writes each user's `NaiveProxy credential` in the format `basic:<base64(username:password)>`.
- Triggers `lua reload users` through HAProxy's UNIX socket for zero-downtime in-memory reload without container restart.
 
### File path
 
- Relative path in module logic: `haproxy/data/users.csv`
- Absolute path inside the node container: `/opt/app/haproxy/data/users.csv`
- Path inside HAProxy container: `/etc/haproxy/data/users.csv`
 
### Row format
 
Strict 2-column format without legacy 1/0 prefixes:
 
`<username>,<credential>`
 
### Example
 
```text
pablo_escobar,90cdd126-819d-5f2f-a1dd-4cbf085182b9
pablo_escobar,basic:cGFibG9fZXNjb2JhcjpteXBhc3N3b3Jk
alvarez,9f37a85cafb66ab9c36f866e5c787b2e42a997716c02337f532a259c
guest,8f49d237686bee189aaedaba1ad434c5580356d8e8b5a3d702bd0228
```
 
### Why it is needed
 
- `haproxy/lua/auth.lua` reads users from `/etc/haproxy/data/users.csv`.
- L4: Trojan is detected by credential hash, including MUX verification.
- L4: VLESS is detected by binary UUID handshake.
- L4: AnyTLS is detected by SHA-256 password hash.
- L7: NaiveProxy is verified via `lua.auth_naive` sample fetch checking `Proxy-Authorization: Basic <token>`.
- `haproxy/lua/users_loader.lua` parses the CSV into lookup tables:
  - `vless[uuid] -> username`
  - `trojan[sha224] -> username`
  - `anytls[sha256] -> username`
  - `naive[token] -> username`
- HAProxy configuration routes authorized connections to corresponding backends or falls back to Nginx.
 
### Important
 
- The module does not modify `haproxy.cfg` or Lua scripts.
- The module manages only `users.csv` and triggers the runtime reload command.
