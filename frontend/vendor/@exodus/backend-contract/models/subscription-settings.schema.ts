import { z } from 'zod';

import { ResponseRulesConfigSchema } from './response-rules';
import { CustomRemarksSchema } from './subscription-settings/custom-remarks.schema';
import { HwidSettingsSchema } from './subscription-settings/hwid-settings.schema';

export const SubscriptionSettingsSchema = z.object({
    uuid: z.uuid(),
    serveJsonAtBaseSubscription: z.boolean(),

    isShowCustomRemarks: z.boolean(),
    customRemarks: CustomRemarksSchema,

    customResponseHeaders: z.nullable(z.record(z.string(), z.string())),

    randomizeHosts: z.boolean(),

    responseRules: z.nullable(ResponseRulesConfigSchema),

    hwidSettings: z.nullable(HwidSettingsSchema),

    createdAt: z.iso.datetime().transform((str) => new Date(str)),
    updatedAt: z.iso.datetime().transform((str) => new Date(str)),
});
