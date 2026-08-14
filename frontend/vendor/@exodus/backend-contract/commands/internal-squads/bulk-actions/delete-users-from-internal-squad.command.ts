import { z } from 'zod';

import { INTERNAL_SQUADS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace DeleteUsersFromInternalSquadCommand {
    export const url = REST_API.INTERNAL_SQUADS.BULK_ACTIONS.REMOVE_USERS;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        INTERNAL_SQUADS_ROUTES.BULK_ACTIONS.REMOVE_USERS(':uuid'),
        'delete',
        'Delete users from internal squad',
        { scope: 'remove-users', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
