import { z } from 'zod';

import { REST_API, SUBSCRIPTIONS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { numberParamSchema } from '../../../models';

export namespace GetConnectionKeysByUserIdCommand {
    export const url = REST_API.SUBSCRIPTIONS.GET_CONNECTION_KEYS_BY_USER_ID;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTIONS_ROUTES.GET_CONNECTION_KEYS_BY_USER_ID(':userId'),
        'get',
        'Get connection keys (base64 format) by user id',
        { scope: 'connection-keys', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        userId: numberParamSchema.describe('User ID'),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            enabledKeys: z.array(z.string()),
            hiddenKeys: z.array(z.string()),
            disabledKeys: z.array(z.string()),
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
