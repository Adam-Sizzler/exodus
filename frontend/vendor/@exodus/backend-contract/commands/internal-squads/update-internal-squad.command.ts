import { z } from 'zod';

import { INTERNAL_SQUADS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { InternalSquadSchema } from '../../models';

export namespace UpdateInternalSquadCommand {
    export const url = REST_API.INTERNAL_SQUADS.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        INTERNAL_SQUADS_ROUTES.UPDATE,
        'patch',
        'Update internal squad',
        { scope: 'update', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuid: z.uuid(),
        name: z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            )
            .optional(),
        inbounds: z.array(z.uuid()).optional(),
    });

    export const ResponseSchema = z.object({
        response: InternalSquadSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
