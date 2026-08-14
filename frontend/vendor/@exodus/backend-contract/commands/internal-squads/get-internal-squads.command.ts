import { z } from 'zod';

import { INTERNAL_SQUADS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { InternalSquadSchema } from '../../models';

export namespace GetInternalSquadsCommand {
    export const url = REST_API.INTERNAL_SQUADS.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        INTERNAL_SQUADS_ROUTES.GET,
        'get',
        'Get all internal squads',
        { scope: 'list', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            total: z.number(),
            internalSquads: z.array(InternalSquadSchema),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
