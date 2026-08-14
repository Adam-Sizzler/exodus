import { z } from 'zod';

import { REST_API, SYSTEM_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace GetRecapCommand {
    export const url = REST_API.SYSTEM.STATS.RECAP;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        SYSTEM_ROUTES.STATS.RECAP,
        'get',
        'Get Recap',
        { scope: 'recap', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            thisMonth: z.object({
                users: z.number(),
                traffic: z.string(),
            }),
            total: z.object({
                users: z.number(),
                nodes: z.number(),
                traffic: z.string(),
                nodesRam: z.string(),
                nodesCpuCores: z.number(),
                distinctCountries: z.number(),
            }),
            version: z.string(),
            initDate: z.iso.datetime({ local: true, offset: true })
                .transform((str) => new Date(str)),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
