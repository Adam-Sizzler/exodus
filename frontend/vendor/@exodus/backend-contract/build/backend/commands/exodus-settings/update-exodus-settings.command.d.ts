import { z } from 'zod';
import {
    BrandingSettingsSchema,
    ExodusSettingsSchema,
    Oauth2SettingsSchema,
    PasskeySettingsSchema,
    PasswordAuthSettingsSchema,
    TgAuthSettingsSchema,
} from '../../models';
export declare namespace UpdateExodusSettingsCommand {
    const url: "/api/exodus-settings/update";
    const TSQ_url: "/api/exodus-settings/update";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        passkeySettings: z.ZodOptional<typeof PasskeySettingsSchema>;
        oauth2Settings: z.ZodOptional<typeof Oauth2SettingsSchema>;
        tgAuthSettings: z.ZodOptional<typeof TgAuthSettingsSchema>;
        passwordSettings: z.ZodOptional<typeof PasswordAuthSettingsSchema>;
        brandingSettings: z.ZodOptional<typeof BrandingSettingsSchema>;
    }, "strip", z.ZodTypeAny, {
        passkeySettings?: z.infer<typeof PasskeySettingsSchema> | undefined;
        oauth2Settings?: z.infer<typeof Oauth2SettingsSchema> | undefined;
        tgAuthSettings?: z.infer<typeof TgAuthSettingsSchema> | undefined;
        passwordSettings?: z.infer<typeof PasswordAuthSettingsSchema> | undefined;
        brandingSettings?: z.infer<typeof BrandingSettingsSchema> | undefined;
    }, {
        passkeySettings?: z.infer<typeof PasskeySettingsSchema> | undefined;
        oauth2Settings?: z.infer<typeof Oauth2SettingsSchema> | undefined;
        tgAuthSettings?: z.infer<typeof TgAuthSettingsSchema> | undefined;
        passwordSettings?: z.infer<typeof PasswordAuthSettingsSchema> | undefined;
        brandingSettings?: z.infer<typeof BrandingSettingsSchema> | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: typeof ExodusSettingsSchema;
    }, "strip", z.ZodTypeAny, {
        response: z.infer<typeof ExodusSettingsSchema>;
    }, {
        response: z.infer<typeof ExodusSettingsSchema>;
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
