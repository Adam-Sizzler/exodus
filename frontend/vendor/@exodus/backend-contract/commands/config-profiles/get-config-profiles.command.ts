import { z } from 'zod';

import { CONFIG_PROFILES_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { ConfigProfileSchema } from '../../models';

export namespace GetConfigProfilesCommand {
    export const url = REST_API.CONFIG_PROFILES.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        CONFIG_PROFILES_ROUTES.GET,
        'get',
        'Get config profiles',
        { scope: 'list', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            total: z.number(),
            configProfiles: z.array(ConfigProfileSchema),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
