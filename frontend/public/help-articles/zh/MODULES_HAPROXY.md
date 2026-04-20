## 模块：HAPROXY

简要说明：该模块会在节点部署时生成并维护 HAProxy Lua auth 使用的凭据表。

### 模块做什么

- 创建并更新**一个**用户文件。
- 写入当前节点的所有用户，不绑定到某个特定 inbound。
- 为每个用户写入 `VLESS UUID`。
- 为每个用户写入 `Trojan credential`，格式为 56 个字符的 SHA-224 hash。
- 将启用状态设置为 `1`。

### 文件路径

- 模块逻辑中的相对路径：`haproxy/data/users.csv`
- 节点容器内的绝对路径：`/app/haproxy/data/users.csv`

### 行格式

`1,username,credential`

### 示例

```text
1,pablo_escobar,90cdd126-819d-5f2f-a1dd-4cbf085182b9
1,alvarez,9f37a85cafb66ab9c36f866e5c787b2e42a997716c02337f532a259c
1,guest,8f49d237686bee189aaedaba1ad434c5580356d8e8b5a3d702bd0228
```

### 为什么需要它

- `haproxy/lua/auth.lua` 会从 `/etc/haproxy/data/users.csv` 读取用户并识别协议。
- Trojan 流量通过 credential 识别，包括 MUX 检查。
- VLESS 流量通过二进制握手中的 UUID 识别。
- `haproxy/lua/users_loader.lua` 会解析 CSV 并生成查找 map。
- `trojan[sha224] -> username`
- `vless[uuid] -> username`
- 随后 `haproxy.cfg` 使用 `lua.identify_protocol`，并将流量路由到 `trojan` 或 `vless` backend。

### 重要

- 该模块不会修改 `haproxy.cfg` 或 Lua 脚本。
- 该模块只管理 `users.csv` 文件。
