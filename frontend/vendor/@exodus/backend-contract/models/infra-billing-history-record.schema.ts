import { z } from 'zod';

import { PartialInfraProviderSchema } from './infra-provider.schema';

export const InfraBillingHistoryRecordSchema = z.object({
    uuid: z.uuid(),
    providerUuid: z.uuid(),
    amount: z.number(),
    billedAt: z.iso.datetime()
        .transform((str) => new Date(str)),
    provider: PartialInfraProviderSchema.omit({
        createdAt: true,
        updatedAt: true,
        loginUrl: true,
    }),
});
