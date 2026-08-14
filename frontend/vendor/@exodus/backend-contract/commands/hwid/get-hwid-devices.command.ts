import { z } from 'zod';

import { HWID_ROUTES, REST_API } from '../../api';
import { getEndpointDetails } from '../../constants';
import { HwidUserDeviceSchema, TanstackQueryRequestQuerySchema } from '../../models';

export namespace GetHwidDevicesCommand {
    export const url = REST_API.HWID.GET_ALL_HWID_DEVICES;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        HWID_ROUTES.GET_ALL_HWID_DEVICES,
        'get',
        'Get HWID devices',
        { scope: 'list', kind: 'read' },
        'Please note that the filters here are primarily intended for use by the frontend and rely on expensive operators such as LIKE under the hood. Misusing these filters may negatively impact the performance of your database.',
    );

    export const RequestQuerySchema = TanstackQueryRequestQuerySchema;

    export const ResponseSchema = z.object({
        response: z.object({
            devices: z.array(HwidUserDeviceSchema),
            total: z.number(),
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
