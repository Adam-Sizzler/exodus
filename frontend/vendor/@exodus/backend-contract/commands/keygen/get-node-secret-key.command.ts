import { z } from 'zod';

import { KEYGEN_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
export namespace GetNodeSecretKeyCommand {
    export const url = REST_API.KEYGEN.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        KEYGEN_ROUTES.GET,
        'get',
        'Get SECRET_KEY for Exodus Node',
        { scope: 'get', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            secretKey: z.string(),
            grpcToken: z.string(),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
