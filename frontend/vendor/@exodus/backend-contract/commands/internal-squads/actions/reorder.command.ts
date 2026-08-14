import { z } from 'zod';

import { INTERNAL_SQUADS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { InternalSquadSchema } from '../../../models';

export namespace ReorderInternalSquadCommand {
    export const url = REST_API.INTERNAL_SQUADS.ACTIONS.REORDER;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        INTERNAL_SQUADS_ROUTES.ACTIONS.REORDER,
        'post',
        'Reorder internal squads',
        { scope: 'reorder', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        items: z.array(
            InternalSquadSchema.pick({
                viewPosition: true,
                uuid: true,
            }),
        ),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            total: z.number(),
            internalSquads: z.array(InternalSquadSchema),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
