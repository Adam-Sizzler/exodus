import { z } from 'zod';

import { REST_API, SUBSCRIPTION_PAGE_CONFIGS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { SubscriptionPageConfigSchema } from '../../models';

export namespace UpdateSubpageConfigCommand {
    export const url = REST_API.SUBSCRIPTION_PAGE_CONFIGS.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_PAGE_CONFIGS_ROUTES.UPDATE,
        'patch',
        'Update subscription page config',
        { scope: 'update', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuid: z.uuid(),
        name: z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            )
            .optional(),
        config: z.optional(z.unknown()),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionPageConfigSchema.extend({
            config: z.unknown(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
