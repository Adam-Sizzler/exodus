## ماژول: HAPROXY

خلاصه: این ماژول جدول اطلاعات احراز هویت مورد استفاده در HAProxy Lua auth را هنگام deploy نود می‌سازد و به‌روز نگه می‌دارد.

### این ماژول چه کاری انجام می‌دهد

- **یک** فایل کاربران را ایجاد و به‌روزرسانی می‌کند.
- همه کاربران نود فعلی را بدون وابستگی به inbound مشخص در فایل قرار می‌دهد.
- برای هر کاربر `VLESS UUID` را می‌نویسد.
- برای هر کاربر `Trojan credential` را به صورت hash نوع SHA-224 با طول 56 کاراکتر می‌نویسد.
- وضعیت فعال بودن را برابر `1` قرار می‌دهد.

### مسیر فایل

- مسیر نسبی در منطق ماژول: `haproxy/data/users.csv`
- مسیر مطلق داخل کانتینر نود: `/app/haproxy/data/users.csv`

### قالب هر خط

`1,username,credential`

### نمونه

```text
1,pablo_escobar,90cdd126-819d-5f2f-a1dd-4cbf085182b9
1,alvarez,9f37a85cafb66ab9c36f866e5c787b2e42a997716c02337f532a259c
1,guest,8f49d237686bee189aaedaba1ad434c5580356d8e8b5a3d702bd0228
```

### چرا لازم است

- `haproxy/lua/auth.lua` کاربران را از `/etc/haproxy/data/users.csv` می‌خواند و پروتکل را تشخیص می‌دهد.
- ترافیک Trojan با credential، همراه با بررسی MUX، تشخیص داده می‌شود.
- ترافیک VLESS با UUID موجود در handshake باینری تشخیص داده می‌شود.
- `haproxy/lua/users_loader.lua` فایل CSV را به mapهای جستجو تبدیل می‌کند.
- `trojan[sha224] -> username`
- `vless[uuid] -> username`
- سپس `haproxy.cfg` از `lua.identify_protocol` استفاده می‌کند و ترافیک را به backendهای `trojan` یا `vless` هدایت می‌کند.

### مهم

- این ماژول `haproxy.cfg` یا اسکریپت‌های Lua را تغییر نمی‌دهد.
- این ماژول فقط فایل `users.csv` را مدیریت می‌کند.
