import { z } from 'zod';

import { REST_API, SUBSCRIPTION_PAGE_CONFIGS_ROUTES } from '../../api';
import { getEndpointDetails } from '../../constants';

export namespace DeleteSubpageConfigCommand {
    export const url = REST_API.SUBSCRIPTION_PAGE_CONFIGS.DELETE;
    export const TSQ_url = url(':uuid');

    export const endpointDetails = getEndpointDetails(
        SUBSCRIPTION_PAGE_CONFIGS_ROUTES.DELETE(':uuid'),
        'delete',
        'Delete subscription page config',
        { scope: 'delete', kind: 'write' },
    );

    export const RequestParamSchema = z.object({
        uuid: z.uuid(),
    });

    export type RequestParam = z.infer<typeof RequestParamSchema>;
}
