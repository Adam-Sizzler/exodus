import { z } from 'zod';

import { NODE_PLUGINS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { NodePluginSchema } from '../../models';

export namespace CreateNodePluginCommand {
    export const url = REST_API.NODE_PLUGINS.CREATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODE_PLUGINS_ROUTES.CREATE,
        'post',
        'Create Node Plugin',
        { scope: 'create', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        name: z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(30, 'Name must be less than 30 characters')
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            ),
    });

    export const ResponseSchema = z.object({
        response: NodePluginSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
