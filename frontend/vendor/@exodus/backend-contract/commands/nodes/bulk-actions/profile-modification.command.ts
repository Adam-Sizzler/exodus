import { z } from 'zod';

import { NODES_ROUTES, REST_API } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace BulkNodesProfileModificationCommand {
    export const url = REST_API.NODES.BULK_ACTIONS.PROFILE_MODIFICATION;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODES_ROUTES.BULK_ACTIONS.PROFILE_MODIFICATION,
        'post',
        'Modify Inbounds & Profile for many nodes',
        { scope: 'bulk-profile-modification', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuids: z.array(z.uuid()).min(1),
        configProfile: z.object({
            activeConfigProfileUuid: z.uuid(),
            activeInbounds: z.array(z.uuid()).min(1),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
