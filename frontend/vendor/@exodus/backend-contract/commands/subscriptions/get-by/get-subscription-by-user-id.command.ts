import { z } from 'zod';

import { REST_API, SUBSCRIPTIONS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { numberParamSchema, SubscriptionInfoSchema } from '../../../models';

export namespace GetSubscriptionByIdCommand {
    export const url = REST_API.SUBSCRIPTIONS.GET_BY.ID;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTIONS_ROUTES.GET_BY.ID(':userId'),
        'get',
        'Get subscription by User ID',
        { scope: 'by-id', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        userId: numberParamSchema.describe('User ID'),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionInfoSchema,
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
