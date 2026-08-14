import { z } from 'zod';

import { REST_API, SUBSCRIPTION_SETTINGS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';
import {
    CustomRemarksSchema,
    HwidSettingsSchema,
    ResponseRulesConfigSchema,
    SubscriptionSettingsSchema,
} from '../../models';

export namespace UpdateSubscriptionSettingsCommand {
    export const url = REST_API.SUBSCRIPTION_SETTINGS.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_SETTINGS_ROUTES.UPDATE,
        'patch',
        'Update subscription settings',
        { scope: 'update', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuid: z.uuid(),
        serveJsonAtBaseSubscription: z.optional(z.boolean()),
        isShowCustomRemarks: z.optional(z.boolean()),
        customRemarks: z.optional(CustomRemarksSchema),
        customResponseHeaders: z.optional(
            z.record(
                z
                    .string()
                    .regex(
                        /^[a-zA-Z0-9_-]+$/,
                        'Invalid header name. Only letters(a-z, A-Z), numbers(0-9), underscores(_) and hyphens(-) are allowed.',
                    ),
                z.string(),
            ),
        ),
        randomizeHosts: z.optional(z.boolean()),
        responseRules: z.optional(ResponseRulesConfigSchema),
        hwidSettings: z.optional(HwidSettingsSchema),
    });

    export const ResponseSchema = z.object({
        response: SubscriptionSettingsSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
