import { z } from 'zod';

import { REST_API, SUBSCRIPTIONS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { SubscriptionInfoSchema } from '../../../models';

export namespace GetSubscriptionByShortUuidProtectedCommand {
    export const url = REST_API.SUBSCRIPTIONS.GET_BY.SHORT_UUID;
    export const TSQ_url = url(':shortUuid');

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTIONS_ROUTES.GET_BY.SHORT_UUID(':shortUuid'),
        'get',
        'Get subscription by short uuid (protected route)',
        { scope: 'by-short-uuid-protected', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        shortUuid: z.string(),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionInfoSchema,
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
