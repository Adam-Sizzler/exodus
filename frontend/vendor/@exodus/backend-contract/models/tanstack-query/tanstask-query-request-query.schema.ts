import { z } from 'zod';

import { TanstackQueryFilterSchema } from './tanstask-query-filter.schema';
import { TanstackQuerySortingSchema } from './tanstask-query-sorting.schema';

export const TanstackQueryRequestQuerySchema = z.object({
    start: z.coerce
        .number()
        .default(0)
        .describe('Start index (offset) of the results to return, default is 0'),
    size: z.coerce
        .number()
        .min(1, 'Size (limit) must be greater than 0')
        .max(1000, 'Size (limit) must be less than 1000')
        .describe('Number of results to return, no more than 1000')
        .default(25),
    filters: z
        .preprocess(
            (str) => (typeof str === 'string' ? JSON.parse(str) : str),
            z.array(TanstackQueryFilterSchema),
        )
        .optional(),

    filterModes: z
        .preprocess(
            (str) => (typeof str === 'string' ? JSON.parse(str) : str),
            z.record(z.string(), z.string()),
        )
        .optional(),

    globalFilterMode: z.string().optional(),

    sorting: z
        .preprocess(
            (str) => (typeof str === 'string' ? JSON.parse(str) : str),
            z.array(TanstackQuerySortingSchema),
        )
        .optional(),
});
