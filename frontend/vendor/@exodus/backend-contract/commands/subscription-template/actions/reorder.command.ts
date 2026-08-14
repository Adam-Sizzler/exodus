import { z } from 'zod';

import { REST_API, SUBSCRIPTION_TEMPLATE_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { SubscriptionTemplateSchema } from '../../../models';

export namespace ReorderSubscriptionTemplateCommand {
    export const url = REST_API.SUBSCRIPTION_TEMPLATE.ACTIONS.REORDER;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_TEMPLATE_ROUTES.ACTIONS.REORDER,
        'post',
        'Reorder subscription templates',
        { scope: 'reorder', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        items: z.array(
            SubscriptionTemplateSchema.pick({
                viewPosition: true,
                uuid: true,
            }),
        ),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            total: z.number(),
            templates: z.array(SubscriptionTemplateSchema),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
