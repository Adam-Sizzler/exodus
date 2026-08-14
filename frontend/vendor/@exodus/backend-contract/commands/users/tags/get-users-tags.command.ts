import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace GetUsersTagsCommand {
    export const url = REST_API.USERS.TAGS.GET;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.TAGS.GET,
        'get',
        'Get users tags',
        { scope: 'list-tags', kind: 'read' },
    );

    export const ResponseSchema = z.object({
        response: z.object({
            tags: z.array(z.string()),
        }),
    });

    export type Response = z.infer<typeof ResponseSchema>;
}
