import z from 'zod';

export const TanstackQueryFilterSchema = z.object({
    id: z.string(),
    value: z.unknown(),
});
