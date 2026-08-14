import { z } from 'zod';

import { REST_API } from '../../api';

export namespace GetSubscriptionByShortUuidCommand {
    export const url = REST_API.SUBSCRIPTION.GET;
    export const TSQ_url = url(':shortUuid');

    export const RequestParamSchema = z.object({
        shortUuid: z.string(),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
