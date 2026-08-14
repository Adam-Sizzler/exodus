import { z } from 'zod';

import { REST_API, SUBSCRIPTION_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { SubscriptionInfoSchema } from '../../models';

export namespace GetSubscriptionInfoByShortUuidCommand {
    export const url = REST_API.SUBSCRIPTION.GET_INFO;
    export const TSQ_url = url(':shortUuid');

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_ROUTES.GET_INFO(':shortUuid'),
        'get',
        'Get Subscription Info by Short UUID',
        { scope: 'get', kind: 'read' },
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
