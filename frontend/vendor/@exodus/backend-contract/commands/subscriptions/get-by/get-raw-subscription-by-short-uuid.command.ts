import { z } from 'zod';

import { REST_API, SUBSCRIPTIONS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { ExtendedUsersSchema, ResolvedProxyConfigSchema } from '../../../models';

export namespace GetRawSubscriptionByShortUuidCommand {
    export const url = REST_API.SUBSCRIPTIONS.GET_BY.SHORT_UUID_RAW;
    export const TSQ_url = url(':shortUuid');

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTIONS_ROUTES.GET_BY.SHORT_UUID_RAW(':shortUuid'),
        'get',
        'Get Raw Subscription by Short UUID',
        { scope: 'raw', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        shortUuid: z.string(),
    });

    export const RequestQuerySchema = z.object({
        withDisabledHosts: z
            .string()
            .transform((str) => str === 'true')
            .optional()
            .prefault('false'),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            user: ExtendedUsersSchema,
            convertedUserInfo: z.object({
                daysLeft: z.number(),
                trafficLimit: z.string(),
                trafficUsed: z.string(),
                lifetimeTrafficUsed: z.string(),
                hwidCheckup: z
                    .object({
                        subscriptionAllowed: z.boolean(),
                        maxDeviceReached: z.boolean(),
                        hwidNotSupported: z.boolean(),
                        limitBypassed: z.boolean(),
                    })
                    .nullable(),
            }),
            headers: z.record(z.string(), z.string().optional()),
            resolvedProxyConfigs: z.array(ResolvedProxyConfigSchema),
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
