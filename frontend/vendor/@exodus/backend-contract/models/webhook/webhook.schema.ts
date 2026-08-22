import z from 'zod';

import { EVENTS, EVENTS_SCOPES, toZodEnum, CRUD_ACTIONS } from '../../constants';
import { ExtendedUsersSchema } from '../extended-users.schema';
import { HwidUserDeviceSchema } from '../hwid-user-device.schema';
import { NodesSchema } from '../nodes.schema';

export const ExodusWebhookUserEvents = z.object({
    scope: z.literal(EVENTS_SCOPES.USER),
    event: z.enum(toZodEnum(EVENTS.USER)),
    timestamp: z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: ExtendedUsersSchema,
    meta: z
        .object({
            notConnectedAfterHours: z.number().nullish(),
            expiration: z.number().nullish(),
        })
        .nullable(),
});

export const ExodusWebhookUserHwidDevicesEvents = z.object({
    scope: z.literal(EVENTS_SCOPES.USER_HWID_DEVICES),
    event: z.enum(toZodEnum(EVENTS.USER_HWID_DEVICES)),
    timestamp: z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: z.object({
        user: ExtendedUsersSchema,
        hwidUserDevice: HwidUserDeviceSchema,
    }),
});

export const ExodusWebhookNodeEvents = z.object({
    scope: z.literal(EVENTS_SCOPES.NODE),
    event: z.enum(toZodEnum(EVENTS.NODE)),
    timestamp: z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: NodesSchema,
});

export const ExodusWebhookServiceEvents = z.object({
    scope: z.literal(EVENTS_SCOPES.SERVICE),
    event: z.enum(toZodEnum(EVENTS.SERVICE)),
    timestamp: z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: z.object({
        loginAttempt: z
            .object({
                username: z.string(),
                ip: z.string(),
                userAgent: z.string(),
                description: z.string().optional(),
                password: z.string().optional(),
            })
            .optional(),
        panelVersion: z.string().optional(),
        subpageConfig: z
            .object({
                action: z.enum(toZodEnum(CRUD_ACTIONS)),
                uuid: z.uuid(),
            })

            .optional(),
        apiToken: z
            .object({
                name: z.string(),
                uuid: z.uuid(),
                expireAt: z
                    .string()
                    .datetime()
                    .transform((str) => new Date(str)),
                scopes: z.array(z.string()),
            })
            .optional(),
    }),
});

export const ExodusWebhookErrorsEvents = z.object({
    scope: z.literal(EVENTS_SCOPES.ERRORS),
    event: z.enum(toZodEnum(EVENTS.ERRORS)),
    timestamp: z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: z.object({
        description: z.string(),
    }),
});

export const ExodusWebhookCrmEvents = z.object({
    scope: z.literal(EVENTS_SCOPES.CRM),
    event: z.enum(toZodEnum(EVENTS.CRM)),
    timestamp: z
        .string()
        .datetime()
        .transform((str) => new Date(str)),
    data: z.object({
        providerName: z.string(),
        nodeName: z.string(),
        nextBillingAt: z
            .string()
            .datetime()
            .transform((str) => new Date(str)),
        loginUrl: z.string(),
    }),
});

export const ExodusWebhookEventSchema = z.discriminatedUnion('scope', [
    ExodusWebhookUserEvents,
    ExodusWebhookUserHwidDevicesEvents,
    ExodusWebhookNodeEvents,
    ExodusWebhookServiceEvents,
    ExodusWebhookErrorsEvents,
    ExodusWebhookCrmEvents,
]);

export type TExodusWebhookEvent = z.infer<typeof ExodusWebhookEventSchema>;

export type TExodusWebhookUserEvent = z.infer<typeof ExodusWebhookUserEvents>;
export type TExodusWebhookNodeEvent = z.infer<typeof ExodusWebhookNodeEvents>;
export type TExodusWebhookServiceEvent = z.infer<typeof ExodusWebhookServiceEvents>;
export type TExodusWebhookErrorsEvent = z.infer<typeof ExodusWebhookErrorsEvents>;
export type TExodusWebhookCrmEvent = z.infer<typeof ExodusWebhookCrmEvents>;
export type TExodusWebhookUserHwidDevicesEvent = z.infer<
    typeof ExodusWebhookUserHwidDevicesEvents
>;
