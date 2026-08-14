import { z } from 'zod';

import { CONFIG_PROFILES_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { ConfigProfileSchema } from '../../models';

export namespace UpdateConfigProfileCommand {
    export const url = REST_API.CONFIG_PROFILES.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        CONFIG_PROFILES_ROUTES.UPDATE,
        'patch',
        'Update Core Config in specific config profile',
        { scope: 'update', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuid: z.uuid('UUID must be a valid UUID'),
        name: z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            )
            .optional(),
        config: z.looseObject({}).optional(),
    });

    export const ResponseSchema = z.object({
        response: ConfigProfileSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
