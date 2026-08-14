import { z } from 'zod';

import { PASSKEYS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace GetPasskeysCommand {
    export const url = REST_API.PASSKEYS.GET_ALL_PASSKEYS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        PASSKEYS_ROUTES.GET_ALL_PASSKEYS,
        'get',
        'Get passkeys',
        { scope: 'list', kind: 'read' },
    );

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

    export type Response = z.infer<typeof ResponseSchema>;
}
