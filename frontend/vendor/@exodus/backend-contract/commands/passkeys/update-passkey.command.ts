import { z } from 'zod';

import { PASSKEYS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace UpdatePasskeyCommand {
    export const url = REST_API.PASSKEYS.UPDATE_PASSKEY;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        PASSKEYS_ROUTES.UPDATE_PASSKEY,
        'patch',
        'Update passkey',
        { scope: 'update', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        id: z.string(),
        name: z
            .string()
            .min(2)
            .max(30)
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            ),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            passkeys: z.array(
                z.object({
                    id: z.string(),
                    name: z.string(),
                    createdAt: z.iso
                        .datetime({ offset: true, local: true })
                        .transform((str) => new Date(str))
                        .describe('Created date. Format: 2025-01-17T15:38:45.065Z'),
                    lastUsedAt: z.iso
                        .datetime({ offset: true, local: true })
                        .transform((str) => new Date(str))
                        .describe('Last used date. Format: 2025-01-17T15:38:45.065Z'),
                }),
            ),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
