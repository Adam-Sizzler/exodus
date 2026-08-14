import { z } from 'zod';

import { HWID_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { HwidUserDeviceSchema } from '../../models';

export namespace CreateUserHwidDeviceCommand {
    export const url = REST_API.HWID.CREATE_USER_HWID_DEVICE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        HWID_ROUTES.CREATE_USER_HWID_DEVICE,
        'post',
        'Create a user HWID device',
        { scope: 'create', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        hwid: z.string().regex(/^[a-zA-Z0-9=-]{10,64}$/),
        userId: z.number(),
        platform: z.optional(z.string()),
        osVersion: z.optional(z.string()),
        deviceModel: z.optional(z.string()),
        userAgent: z.optional(z.string()),
        requestIp: z.optional(z.string()),
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
