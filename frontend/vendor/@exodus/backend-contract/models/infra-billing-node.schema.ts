import { z } from 'zod';

import { PartialInfraProviderSchema } from './infra-provider.schema';
import { NodesSchema } from './nodes.schema';

export const InfraBillingNodeSchema = z.object({
    uuid: z.uuid(),
    nodeUuid: z.uuid().nullable(),
    name: z.string().nullable(),
    providerUuid: z.uuid(),
    provider: PartialInfraProviderSchema.pick({
        uuid: true,
        name: true,
        loginUrl: true,
        faviconLink: true,
    }),

    node: NodesSchema.pick({
        uuid: true,
        name: true,
        countryCode: true,
    }).nullable(),

    nextBillingAt: z.iso.datetime()
        .transform((str) => new Date(str)),

    createdAt: z.iso.datetime()
        .transform((str) => new Date(str)),
    updatedAt: z.iso.datetime()
        .transform((str) => new Date(str)),
});
