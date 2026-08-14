import { z } from 'zod';

import { HWID_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { HwidUserDeviceSchema } from '../../models';

export namespace DeleteUserHwidDeviceCommand {
    export const url = REST_API.HWID.DELETE_USER_HWID_DEVICE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        HWID_ROUTES.DELETE_USER_HWID_DEVICE,
        'post',
        'Delete a user HWID device',
        { scope: 'delete', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        userId: z.number(),
        hwid: z.string(),
    });

    export const ResponseSchema = z.object({
        response: z.object({
            total: z.number(),
            devices: z.array(HwidUserDeviceSchema),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
