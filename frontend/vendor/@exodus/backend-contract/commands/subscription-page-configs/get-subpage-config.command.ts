import { z } from 'zod';

import { REST_API, SUBSCRIPTION_PAGE_CONFIGS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { SubscriptionPageConfigSchema } from '../../models';

export namespace GetSubpageConfigCommand {
    export const url = REST_API.SUBSCRIPTION_PAGE_CONFIGS.GET;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_PAGE_CONFIGS_ROUTES.GET(':uuid'),
        'get',
        'Get subscription page config by uuid',
        { scope: 'get', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionPageConfigSchema.extend({
            config: z.unknown(),
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
