import { z } from 'zod';

import { EXODUS_SETTINGS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import {
    BrandingSettingsSchema,
    ExodusSettingsSchema,
    Oauth2SettingsSchema,
    PasskeySettingsSchema,
    PasswordAuthSettingsSchema,
} from '../../models';

export namespace UpdateExodusSettingsCommand {
    export const url = REST_API.EXODUS_SETTINGS.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        EXODUS_SETTINGS_ROUTES.UPDATE,
        'patch',
        'Update Exodus settings',
        { scope: 'update', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        passkeySettings: PasskeySettingsSchema.optional(),
        oauth2Settings: Oauth2SettingsSchema.optional(),
        passwordSettings: PasswordAuthSettingsSchema.optional(),
        brandingSettings: BrandingSettingsSchema.optional(),
    });

    export const ResponseSchema = z.object({
        response: ExodusSettingsSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
