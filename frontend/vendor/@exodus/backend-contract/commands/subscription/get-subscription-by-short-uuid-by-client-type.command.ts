import { z } from 'zod';

import { REST_API } from '../../api';
import { REQUEST_TEMPLATE_TYPE } from '../../constants';

export namespace GetSubscriptionByShortUuidByClientTypeCommand {
    export const url = REST_API.SUBSCRIPTION.GET;
    export const TSQ_url = url(':shortUuid');

    export const RequestParamSchema = z.object({
        shortUuid: z.string(),
        clientType: z.enum(REQUEST_TEMPLATE_TYPE, {
            error: 'Invalid client type.'
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
