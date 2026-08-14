import { z } from 'zod';

import { REST_API, SUBSCRIPTION_REQUEST_HISTORY_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { SubscriptionRequestHistorySchema, TanstackQueryRequestQuerySchema } from '../../models';

export namespace GetSubscriptionRequestHistoryCommand {
    export const url = REST_API.SUBSCRIPTION_REQUEST_HISTORY.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_REQUEST_HISTORY_ROUTES.GET,
        'get',
        'Get all subscription request history',
        { scope: 'list', kind: 'read' },
        'Please note that the filters here are primarily intended for use by the frontend and rely on expensive operators such as LIKE under the hood. Misusing these filters may negatively impact the performance of your database.',
    );

    export const RequestQuerySchema = TanstackQueryRequestQuerySchema;

    export const ResponseSchema = z.object({
        response: z.object({
            records: z.array(SubscriptionRequestHistorySchema),
            total: z.number(),
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
