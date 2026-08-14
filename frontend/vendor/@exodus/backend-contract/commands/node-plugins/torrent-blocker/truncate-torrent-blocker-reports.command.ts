import { REST_API, NODE_PLUGINS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';

export namespace TruncateTorrentBlockerReportsCommand {
    export const url = REST_API.NODE_PLUGINS.TORRENT_BLOCKER.TRUNCATE_REPORTS;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        NODE_PLUGINS_ROUTES.TORRENT_BLOCKER.TRUNCATE_REPORTS,
        'delete',
        'Truncate Torrent Blocker Reports',
        { scope: 'truncate', kind: 'write' },
    );
}
