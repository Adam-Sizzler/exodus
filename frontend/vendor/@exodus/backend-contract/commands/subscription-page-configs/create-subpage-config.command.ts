import { z } from 'zod';

import { REST_API, SUBSCRIPTION_PAGE_CONFIGS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { SubscriptionPageConfigSchema } from '../../models';

export namespace CreateSubpageConfigCommand {
    export const url = REST_API.SUBSCRIPTION_PAGE_CONFIGS.CREATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_PAGE_CONFIGS_ROUTES.CREATE,
        'post',
        'Create subscription page config',
        { scope: 'create', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        name: z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            ),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionPageConfigSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
