import { z } from 'zod';

export const BaseStatSchema = z.object({
    current: z.string(),
    previous: z.string(),
    difference: z.string(),
});
