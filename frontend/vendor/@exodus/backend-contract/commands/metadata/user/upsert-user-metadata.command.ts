import { z } from 'zod';

import { REST_API, METADATA_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace UpsertUserMetadataCommand {
    export const url = REST_API.METADATA.USER.UPSERT;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        METADATA_ROUTES.USER.UPSERT(':userId'),
        'put',
        'Update or create User Metadata',
        { scope: 'upsert-user', kind: 'write' },
    );

    export const RequestParamsSchema = z.object({
        userId: z.coerce.number(),
    });

    export const RequestBodySchema = z.object({
        metadata: z.looseObject({}),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            metadata: z.looseObject({}),
        }),
    });

    export type RequestParams = z.infer<typeof RequestParamsSchema>;
    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
