import { z } from 'zod';

import { INTERNAL_SQUADS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace AddManyUsersToInternalSquadCommand {
    export const url = REST_API.INTERNAL_SQUADS.BULK_ACTIONS.ADD_MANY_USERS;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        INTERNAL_SQUADS_ROUTES.BULK_ACTIONS.ADD_MANY_USERS(':uuid'),
        'post',
        'Add many users to internal squad',
        { scope: 'add-many-users', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export const RequestBodySchema = z.object({
        userIds: z.array(z.number()).min(1).max(1000),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
