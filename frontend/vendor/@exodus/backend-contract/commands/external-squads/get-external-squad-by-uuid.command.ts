import { z } from 'zod';

import { EXTERNAL_SQUADS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { ExternalSquadSchema } from '../../models';

export namespace GetExternalSquadByUuidCommand {
    export const url = REST_API.EXTERNAL_SQUADS.GET_BY_UUID;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        EXTERNAL_SQUADS_ROUTES.GET_BY_UUID(':uuid'),
        'get',
        'Get external squad by uuid',
        { scope: 'get', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid().describe('UUID of the external squad'),
    });

    export const ResponseSchema = z.object({
        response: ExternalSquadSchema,
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
