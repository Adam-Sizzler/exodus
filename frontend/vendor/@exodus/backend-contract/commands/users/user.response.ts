import z from 'zod';

import { ExtendedUsersSchema } from '../../models/extended-users.schema';

export const UserResponseSchema = z.object({
    response: ExtendedUsersSchema,
});
