import { z } from 'zod';

import { REST_API, METADATA_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace GetNodeMetadataCommand {
    export const url = REST_API.METADATA.NODE.GET;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        METADATA_ROUTES.NODE.GET(':uuid'),
        'get',
        'Get node metadata',
        { scope: 'get-node', kind: 'read' },
    );

    export const RequestParamsSchema = z.object({
        uuid: z.uuid(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            metadata: z.looseObject({}),
        }),
    });

    export type RequestParams = z.infer<typeof RequestParamsSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
