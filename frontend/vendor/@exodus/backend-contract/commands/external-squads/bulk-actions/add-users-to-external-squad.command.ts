import { z } from 'zod';

import { EXTERNAL_SQUADS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace AddUsersToExternalSquadCommand {
    export const url = REST_API.EXTERNAL_SQUADS.BULK_ACTIONS.ADD_USERS;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        EXTERNAL_SQUADS_ROUTES.BULK_ACTIONS.ADD_USERS(':uuid'),
        'post',
        'Add all users to external squad',
        { scope: 'add-users', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid().describe('UUID of the external squad'),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
