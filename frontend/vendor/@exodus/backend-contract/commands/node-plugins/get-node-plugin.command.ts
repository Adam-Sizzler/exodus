import { z } from 'zod';

import { NODE_PLUGINS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { NodePluginSchema } from '../../models';

export namespace GetNodePluginCommand {
    export const url = REST_API.NODE_PLUGINS.GET;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        NODE_PLUGINS_ROUTES.GET(':uuid'),
        'get',
        'Get Node Plugin by uuid',
        { scope: 'get', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export const ResponseSchema = z.object({
        response: NodePluginSchema.extend({
            pluginConfig: z.unknown(),
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
