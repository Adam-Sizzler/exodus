## Telegram OAuth2

> 本指南适用于 **2.7.0** 及以上版本。

要配置"通过 Telegram 登录"功能，您需要一个 Telegram 机器人。此外，还需要对机器人进行额外配置，功能才能正常工作。

### 机器人配置

1. 打开 @BotFather（https://t.me/botfather）
2. 点击"Open"按钮打开 MiniApp
3. 选择您的机器人，然后点击 `Bot Settings`
4. 打开 `Web Login` 部分。
5. 使用 `/setdomain` 或 `Web Login` 中的域名字段设置面板域名。
   只填写域名，不要填写协议和路径：

- `panel.domain.com`

如果面板地址是 `https://panel.domain.com/panel/auth/login`，BotFather 中应填写 `panel.domain.com`。

### 访问权限配置

填写 `Bot Token` 后，您需要指定允许登录的管理员 ID 列表。

1. 使用对应账号打开机器人 – https://t.me/Get_myidrobot
2. 机器人会返回您的用户 ID，将该 ID 填入对应的字段中。

---

### 已知错误解决方案

###### 安装在面板之上的各种保护措施（如 cookie 等）可能无法与 `Telegram OAuth2` 正常配合工作。

请在反向代理中使用 /oauth2/ 路径来解决此问题

###### 错误：BOT_DOMAIN_INVALID

此错误发生在 backend callback 之前，Telegram 会比较页面 `origin` 和机器人域名。对于 `https://panel.domain.com/panel/auth/login`，Telegram 检查的是 `https://panel.domain.com`，因此 BotFather 中必须填写 `panel.domain.com`，不要包含 `/panel/auth/login`。

###### 错误：登录时未收到 Telegram 验证码

很遗憾，此问题无法从 Exodus 端解决。请尝试使用较早创建的机器人，或使用其他浏览器。
此外，您也可以尝试先登录"官方"资源 — 例如 https://fragment.com。由于浏览器中的 Telegram 会话是共享的 — 之后您可以尝试登录面板。
