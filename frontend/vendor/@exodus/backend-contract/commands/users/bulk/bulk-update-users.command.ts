import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../../api';
import { getEndpointDetails } from '../../../constants';
import { RESET_PERIODS } from '../../../constants';
import { UsersSchema } from '../../../models';

export namespace BulkUpdateUsersCommand {
    export const url = REST_API.USERS.BULK.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.BULK.UPDATE,
        'post',
        'Bulk update users by User IDs',
        { scope: 'bulk-update-users', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        userIds: z.array(z.number()).min(1).max(500),
        fields: z.object({
            status: UsersSchema.shape.status.optional(),
            trafficLimitBytes: z.optional(
                z.number().min(0).describe('Traffic limit in bytes. 0 - unlimited'),
            ),
            trafficLimitStrategy: z.optional(
                z.enum(RESET_PERIODS).describe('Available reset periods'),
            ),
            expireAt: z.optional(
                z.iso
                    .datetime({ local: true, offset: true })
                    .transform((str) => new Date(str))
                    .refine((date) => date > new Date(), {
                        error: 'Expiration date cannot be in the past',
                    })
                    .describe('Expiration date: 2025-01-17T15:38:45.065Z'),
            ),
            description: z.string().nullish(),
            telegramId: z.number().nullish(),
            email: z.email().nullish(),
            tag: z.optional(
                z
                    .string()
                    .regex(
                        /^[A-Z0-9_]+$/,
                        'Tag can only contain uppercase letters, numbers, underscores',
                    )
                    .max(16, 'Tag must be less than 16 characters')
                    .nullable(),
            ),
            hwidDeviceLimit: z.int().min(0).nullish(),
            externalSquadUuid: z.uuid().nullish().describe('Optional. External squad UUID.'),
        }),
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
}
