import { z } from 'zod';

import { INTERNAL_SQUADS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeleteInternalSquadCommand {
    export const url = REST_API.INTERNAL_SQUADS.DELETE;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        INTERNAL_SQUADS_ROUTES.DELETE(':uuid'),
        'delete',
        'Delete internal squad',
        { scope: 'delete', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
