import { z } from 'zod';

export const InfraProviderSchema = z.object({
    uuid: z.uuid(),
    name: z.string(),
    faviconLink: z.nullable(z.string()),
    loginUrl: z.nullable(z.string()),

    createdAt: z.iso.datetime()
        .transform((str) => new Date(str)),
    updatedAt: z.iso.datetime()
        .transform((str) => new Date(str)),

    billingHistory: z.object({
        totalAmount: z.number(),
        totalBills: z.number(),
    }),
    billingNodes: z.array(
        z.object({
            name: z.string(),
            details: z
                .object({
                    nodeUuid: z.uuid(),
                    countryCode: z.string(),
                })
                .nullable(),
        }),
    ),
});

export const PartialInfraProviderSchema = z.object({
    uuid: z.uuid(),
    name: z.string(),
    faviconLink: z.nullable(z.string()),
    loginUrl: z.nullable(z.string()),

    createdAt: z.iso.datetime()
        .transform((str) => new Date(str)),
    updatedAt: z.iso.datetime()
        .transform((str) => new Date(str)),
});
