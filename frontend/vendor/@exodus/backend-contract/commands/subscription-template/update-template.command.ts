import { z } from 'zod';

import { REST_API, SUBSCRIPTION_TEMPLATE_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { SubscriptionTemplateSchema } from '../../models';

export namespace UpdateSubscriptionTemplateCommand {
    export const url = REST_API.SUBSCRIPTION_TEMPLATE.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_TEMPLATE_ROUTES.UPDATE,
        'patch',
        'Update subscription template',
        { scope: 'update', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuid: z.uuid(),
        name: z
            .string()
            .min(2)
            .max(255)
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            )
            .optional(),
        templateJson: z.optional(z.looseObject({})),
        encodedTemplateYaml: z.optional(z.string()),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionTemplateSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
