import { z } from 'zod';

import { REST_API, USERS_ROUTES } from '../../api';
import { getEndpointDetails, RESET_PERIODS, USERS_STATUS } from '../../constants';
import { UserResponseSchema } from './user.response';
export namespace UpdateUserCommand {
    export const url = REST_API.USERS.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        USERS_ROUTES.UPDATE,
        'patch',
        'Update a user',
        { scope: 'update', kind: 'write' },
        'Update a user by ID or username. Exactly one of the fields must be provided.',
    );

    export const RequestBodySchema = z
        .object({
            username: z.optional(z.string().describe('Username of the user')),
            id: z.optional(z.number().describe('ID of the user')),
            status: z.enum([USERS_STATUS.ACTIVE, USERS_STATUS.DISABLED]).optional(),
            trafficLimitBytes: z
                .number()
                .min(0)
                .describe('Traffic limit in bytes. 0 - unlimited')
                .optional(),
            trafficLimitStrategy: z
                .enum(RESET_PERIODS)
                .describe('Traffic limit reset strategy')
                .optional(),
            expireAt: z.iso
                .datetime({ local: true, offset: true })
                .transform((str) => new Date(str))
                .refine((date) => date > new Date(), {
                    error: 'Expiration date cannot be in the past',
                })
                .describe('Expiration date: 2025-01-17T15:38:45.065Z')
                .optional(),
            description: z.optional(z.string().nullable()),
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
            telegramId: z.number().nullish(),
            email: z.email().nullish(),
            hwidDeviceLimit: z.int().min(0).nullish(),
            trojanPassword: z.string().min(8).max(32).optional(),
            vlessUuid: z.uuid().optional(),
            ssPassword: z.string().min(8).max(32).optional(),
            naivePassword: z.string().min(8).max(32).optional(),
            shadowtlsPassword: z.string().min(8).max(32).optional(),
            hysteria2Password: z.string().min(8).max(32).optional(),
            anytlsPassword: z.string().min(8).max(32).optional(),
            activeInternalSquads: z.array(z.uuid()).optional(),
            externalSquadUuid: z
                .optional(z.nullable(z.uuid()))
                .describe('Optional. External squad UUID.'),
        })
        .refine((d) => d.username ?? d.id, {
            error: 'At least one of username, id must be provided',
        });

    export const ResponseSchema = UserResponseSchema;

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
