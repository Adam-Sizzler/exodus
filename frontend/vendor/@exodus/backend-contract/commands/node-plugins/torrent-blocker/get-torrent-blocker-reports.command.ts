import { z } from 'zod';

import { REST_API, NODE_PLUGINS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { TorrentBlockerReportSchema, TanstackQueryRequestQuerySchema } from '../../../models';

export namespace GetTorrentBlockerReportsCommand {
    export const url = REST_API.NODE_PLUGINS.TORRENT_BLOCKER.GET_REPORTS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODE_PLUGINS_ROUTES.TORRENT_BLOCKER.GET_REPORTS,
        'get',
        'Get Torrent Blocker Reports',
        { scope: 'torrent-blocker-reports', kind: 'read' },
        'Please note that the filters here are primarily intended for use by the frontend and rely on expensive operators such as LIKE under the hood. Misusing these filters may negatively impact the performance of your database.',
    );

    export const RequestQuerySchema = TanstackQueryRequestQuerySchema;

    export const ResponseSchema = z.object({
        response: z.object({
            records: z.array(TorrentBlockerReportSchema),
            total: z.number(),
        }),
    });

    export type RequestQuery = z.infer<typeof RequestQuerySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
