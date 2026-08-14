import { z } from 'zod';

import { CONFIG_PROFILES_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeleteConfigProfileCommand {
    export const url = REST_API.CONFIG_PROFILES.DELETE;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        CONFIG_PROFILES_ROUTES.DELETE(':uuid'),
        'delete',
        'Delete config profile',
        { scope: 'delete', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid().describe('UUID of the config profile'),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
