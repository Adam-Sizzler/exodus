import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { UserResponseSchema } from '../user.response';

export namespace GetUserByShortUuidCommand {
    export const url = REST_API.USERS.GET_BY.SHORT_UUID;
    export const TSQ_url = url(':shortUuid');

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.GET_BY.SHORT_UUID(':shortUuid'),
        'get',
        'Get user by Short UUID',
        { scope: 'by-short-uuid', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        shortUuid: z.string(),
    });

    export const ResponseSchema = UserResponseSchema;

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
