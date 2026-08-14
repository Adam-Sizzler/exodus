import { z } from 'zod';

import { REST_API, METADATA_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { numberParamSchema } from '../../../models';

export namespace GetUserMetadataCommand {
    export const url = REST_API.METADATA.USER.GET;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        METADATA_ROUTES.USER.GET(':userId'),
        'get',
        'Get user metadata',
        { scope: 'get-user', kind: 'read' },
    );

    export const RequestParamsSchema = z.object({
        userId: numberParamSchema,
    });

    export const ResponseSchema = z.object({
        response: z.object({
            metadata: z.looseObject({}),
        }),
    });

    export type RequestParams = z.infer<typeof RequestParamsSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
