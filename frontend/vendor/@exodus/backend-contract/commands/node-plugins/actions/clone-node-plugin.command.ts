import { z } from 'zod';

import { NODE_PLUGINS_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { NodePluginSchema } from '../../../models';

export namespace CloneNodePluginCommand {
    export const url = REST_API.NODE_PLUGINS.ACTIONS.CLONE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODE_PLUGINS_ROUTES.ACTIONS.CLONE,
        'post',
        'Clone Node Plugin',
        { scope: 'clone', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        cloneFromUuid: z.uuid(),
    });

    export const ResponseSchema = z.object({
        response: NodePluginSchema.extend({
            pluginConfig: z.unknown(),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
