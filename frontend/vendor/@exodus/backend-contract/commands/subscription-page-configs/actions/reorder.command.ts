import { z } from 'zod';

import { REST_API, SUBSCRIPTION_PAGE_CONFIGS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { SubscriptionPageConfigSchema } from '../../../models';

export namespace ReorderSubpageConfigsCommand {
    export const url = REST_API.SUBSCRIPTION_PAGE_CONFIGS.ACTIONS.REORDER;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_PAGE_CONFIGS_ROUTES.ACTIONS.REORDER,
        'post',
        'Reorder subscription page configs',
        { scope: 'reorder', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        items: z.array(
            SubscriptionPageConfigSchema.pick({
                viewPosition: true,
                uuid: true,
            }),
        ),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            total: z.number(),
            configs: z.array(SubscriptionPageConfigSchema),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
