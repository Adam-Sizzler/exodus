import { z } from 'zod';

import { REST_API, SUBSCRIPTION_TEMPLATE_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeleteSubscriptionTemplateCommand {
    export const url = REST_API.SUBSCRIPTION_TEMPLATE.DELETE;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_TEMPLATE_ROUTES.DELETE(':uuid'),
        'delete',
        'Delete subscription template',
        { scope: 'delete', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
