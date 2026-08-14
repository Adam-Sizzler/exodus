import { z } from 'zod';

import { REST_API, SUBSCRIPTION_TEMPLATE_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { SubscriptionTemplateSchema } from '../../models';

export namespace GetSubscriptionTemplatesCommand {
    export const url = REST_API.SUBSCRIPTION_TEMPLATE.GET_ALL;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_TEMPLATE_ROUTES.GET_ALL,
        'get',
        'Get all subscription templates (wihout content)',
        { scope: 'list', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            total: z.number(),
            templates: z.array(SubscriptionTemplateSchema),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
