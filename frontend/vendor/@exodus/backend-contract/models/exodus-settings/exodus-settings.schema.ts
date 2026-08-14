import { z } from 'zod';

import { BrandingSettingsSchema } from './branding-settings.schema';
import { Oauth2SettingsSchema } from './oauth2-settings.schema';
import { PasskeySettingsSchema } from './passkey-settings.schema';
import { PasswordAuthSettingsSchema } from './password-auth-settings.schema';

export const ExodusSettingsSchema = z.object({
    passkeySettings: z.nullable(PasskeySettingsSchema),
    oauth2Settings: z.nullable(Oauth2SettingsSchema),
    passwordSettings: z.nullable(PasswordAuthSettingsSchema),
    brandingSettings: z.nullable(BrandingSettingsSchema),
});

export type TExodusSettings = z.infer<typeof ExodusSettingsSchema>;
