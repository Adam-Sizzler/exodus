import { z } from 'zod';

export const BaseInternalSquadSchema = z.object({
    uuid: z.uuid(),
    name: z.string(),
});
