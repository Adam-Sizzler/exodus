import { z } from 'zod';

import { HostsSchema } from '../../models/hosts.schema';

export const HostResponseSchema = z.object({
    response: HostsSchema,
});
