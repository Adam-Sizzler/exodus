import { z } from 'zod';

import { REST_API, SUBSCRIPTION_PAGE_CONFIGS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { SubscriptionPageConfigSchema } from '../../../models';

export namespace CloneSubpageConfigCommand {
    export const url = REST_API.SUBSCRIPTION_PAGE_CONFIGS.ACTIONS.CLONE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_PAGE_CONFIGS_ROUTES.ACTIONS.CLONE,
        'post',
        'Clone subscription page config',
        { scope: 'clone', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        cloneFromUuid: z.uuid(),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionPageConfigSchema.extend({
            config: z.unknown(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
