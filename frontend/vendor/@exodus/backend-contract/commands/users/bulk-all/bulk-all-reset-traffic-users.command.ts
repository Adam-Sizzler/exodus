import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkAllResetTrafficUsersCommand {
    export const url = REST_API.USERS.BULK.ALL.RESET_TRAFFIC;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.BULK.ALL.RESET_TRAFFIC,
        'post',
        'Reset user used traffic for all users',
        { scope: 'bulk-all-reset-traffic', kind: 'write' },
    );
}
