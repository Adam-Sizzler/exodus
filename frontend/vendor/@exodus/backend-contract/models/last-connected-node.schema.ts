import { z } from 'zod';

export const LastConnectedNodeSchema = z
    .object({
        connectedAt: z.iso.datetime()
            .transform((str) => new Date(str)),
        nodeName: z.string(),
        countryCode: z.string(),
    })
    .nullable();
