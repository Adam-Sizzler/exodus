import { z } from 'zod';

import { HWID_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { HwidUserDeviceSchema } from '../../models';

export namespace GetUserHwidDevicesCommand {
    export const url = REST_API.HWID.GET_USER_HWID_DEVICES;
    export const TSQ_url = url(':userId');

    export const endpointDetails = getEndpointDetails(
        HWID_ROUTES.GET_USER_HWID_DEVICES(':userId'),
        'get',
        'Get user HWID devices',
        { scope: 'list-by-user', kind: 'read' },
    );

    export const RequestParamSchema = z.object({
        userId: z.coerce.number(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            total: z.number(),
            devices: z.array(HwidUserDeviceSchema),
        }),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
