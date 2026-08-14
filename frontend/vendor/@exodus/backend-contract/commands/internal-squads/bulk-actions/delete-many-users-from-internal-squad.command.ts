import { z } from 'zod';

import { INTERNAL_SQUADS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace DeleteManyUsersFromInternalSquadCommand {
    export const url = REST_API.INTERNAL_SQUADS.BULK_ACTIONS.REMOVE_MANY_USERS;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        INTERNAL_SQUADS_ROUTES.BULK_ACTIONS.REMOVE_MANY_USERS(':uuid'),
        'delete',
        'Delete many users from internal squad',
        { scope: 'remove-many-users', kind: 'write' },
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
