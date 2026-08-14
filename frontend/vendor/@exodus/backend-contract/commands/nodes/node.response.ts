import { z } from 'zod';

import { NodesSchema } from '../../models';

export const NodeResponseSchema = z.object({
    response: NodesSchema,
});
