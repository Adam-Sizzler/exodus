import { z } from 'zod';

import { REST_API, SUBSCRIPTION_TEMPLATE_ROUTES } from '../../api';
import { SUBSCRIPTION_TEMPLATE_TYPE, getEndpointDetails } from '../../constants';
import { SubscriptionTemplateSchema } from '../../models';

export namespace CreateSubscriptionTemplateCommand {
    export const url = REST_API.SUBSCRIPTION_TEMPLATE.CREATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_TEMPLATE_ROUTES.CREATE,
        'post',
        'Create subscription template',
        { scope: 'create', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        name: z
            .string()
            .min(2)
            .max(255)
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            ),
        templateType: z.enum(SUBSCRIPTION_TEMPLATE_TYPE),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionTemplateSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
