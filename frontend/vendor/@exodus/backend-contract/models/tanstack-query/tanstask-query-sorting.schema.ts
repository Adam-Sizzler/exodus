import { z } from 'zod';

export const TanstackQuerySortingSchema = z.object({
    id: z.string(),
    desc: z.boolean(),
});
