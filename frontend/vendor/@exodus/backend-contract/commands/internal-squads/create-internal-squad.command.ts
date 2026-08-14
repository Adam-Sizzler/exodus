import { z } from 'zod';

import { INTERNAL_SQUADS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { InternalSquadSchema } from '../../models';

export namespace CreateInternalSquadCommand {
    export const url = REST_API.INTERNAL_SQUADS.CREATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        INTERNAL_SQUADS_ROUTES.CREATE,
        'post',
        'Create internal squad',
        { scope: 'create', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        name: z
            .string()
            .min(2)
            .max(30)
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            ),
        inbounds: z.array(z.uuid()),
    });

    export const ResponseSchema = z.object({
        response: InternalSquadSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
