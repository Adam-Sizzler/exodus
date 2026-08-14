import { z } from 'zod';

import { NODE_PLUGINS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { NodePluginSchema } from '../../models';

export namespace UpdateNodePluginCommand {
    export const url = REST_API.NODE_PLUGINS.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODE_PLUGINS_ROUTES.UPDATE,
        'patch',
        'Update Node Plugin',
        { scope: 'update', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuid: z.uuid(),
        name: z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            )
            .optional(),
        pluginConfig: z.optional(z.unknown()),
    });

    export const ResponseSchema = z.object({
        response: NodePluginSchema.extend({
            pluginConfig: z.unknown(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
