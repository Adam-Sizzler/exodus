import { z } from 'zod';

import { REST_API, SUBSCRIPTION_SETTINGS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import { SubscriptionSettingsSchema } from '../../models';

export namespace GetSubscriptionSettingsCommand {
    export const url = REST_API.SUBSCRIPTION_SETTINGS.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_SETTINGS_ROUTES.GET,
        'get',
        'Get subscription settings',
        { scope: 'get', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: SubscriptionSettingsSchema,
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
