## Telegram OAuth2

> This guide is applicable for version **2.7.0** and above.

To configure the "Login via Telegram" feature, you need a Telegram bot. Additionally, you need to configure the bot for the feature to work correctly.

### Bot Configuration

1. Open @BotFather (https://t.me/botfather)
2. Open the MiniApp by pressing "Open"
3. Select your bot and press `Bot Settings`
4. Open the `Web Login` section.
5. Set the panel domain with `/setdomain` or the domain field in `Web Login`.
   Specify only the domain, without scheme and path:

- `panel.domain.com`

For a panel opened at `https://data.s-backup.online/panel/auth/login`, BotFather must contain `data.s-backup.online`.

### Access Configuration

After filling in `Bot Token`, you need to specify a list of administrator IDs who will have access to login.

1. From the required account, launch the bot – https://t.me/Get_myidrobot
2. In response, the bot will send you your ID, enter it in the corresponding field.

---

### Known Error Solutions

###### Various protections installed on top of the panel (such as cookies, etc.) may not work correctly with `Telegram OAuth2`.

Use the /oauth2/ path in your reverse proxy to resolve this issue

###### Error: BOT_DOMAIN_INVALID

This error happens before the backend callback when Telegram compares the page `origin` with the bot domain. For `https://data.s-backup.online/panel/auth/login`, Telegram checks `https://data.s-backup.online`, so BotFather must contain `data.s-backup.online`, without `/panel/auth/login`.

###### Error: Telegram confirmation code not received during login

Unfortunately, this issue cannot be resolved from the Exodus side. Try using a bot that was created a while ago or use a different browser.
Alternatively, you can try logging in on one of the "official" resources – for example, https://fragment.com. Since the Telegram session within the browser will be shared – you can try logging into the panel.
