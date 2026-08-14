import { z } from 'zod';

import { REST_API, SUBSCRIPTIONS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { SubscriptionInfoSchema } from '../../../models';

export namespace GetSubscriptionByUsernameCommand {
    export const url = REST_API.SUBSCRIPTIONS.GET_BY.USERNAME;
    export const TSQ_url = url(':username');

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTIONS_ROUTES.GET_BY.USERNAME(':username'),
        'get',
        'Get subscription by username',
        { scope: 'by-username', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        username: z.string().describe('Username'),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionInfoSchema,
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
