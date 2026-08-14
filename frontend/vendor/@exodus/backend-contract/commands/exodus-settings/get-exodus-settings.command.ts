import { z } from 'zod';

import { EXODUS_SETTINGS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { ExodusSettingsSchema } from '../../models';

export namespace GetExodusSettingsCommand {
    export const url = REST_API.EXODUS_SETTINGS.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        EXODUS_SETTINGS_ROUTES.GET,
        'get',
        'Get Exodus settings',
        { scope: 'get', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: ExodusSettingsSchema,
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
