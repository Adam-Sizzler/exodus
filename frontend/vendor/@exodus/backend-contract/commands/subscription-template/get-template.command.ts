import { z } from 'zod';

import { REST_API, SUBSCRIPTION_TEMPLATE_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { SubscriptionTemplateSchema } from '../../models';

export namespace GetSubscriptionTemplateCommand {
    export const url = REST_API.SUBSCRIPTION_TEMPLATE.GET;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_TEMPLATE_ROUTES.GET(':uuid'),
        'get',
        'Get subscription template by uuid',
        { scope: 'get', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionTemplateSchema,
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
