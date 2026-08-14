import { z } from 'zod';

import { EXTERNAL_SQUADS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeleteExternalSquadCommand {
    export const url = REST_API.EXTERNAL_SQUADS.DELETE;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        EXTERNAL_SQUADS_ROUTES.DELETE(':uuid'),
        'delete',
        'Delete external squad',
        { scope: 'delete', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid().describe('UUID of the external squad'),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
