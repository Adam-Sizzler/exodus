import { z } from 'zod';
export declare namespace UpdateUserCommand {
    const url: "/api/users/";
    const TSQ_url: "/api/users/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodEffects<z.ZodObject<{
        username: z.ZodOptional<z.ZodString>;
        uuid: z.ZodOptional<z.ZodString>;
        status: z.ZodOptional<z.ZodEnum<["ACTIVE", "DISABLED"]>>;
        trafficLimitBytes: z.ZodOptional<z.ZodNumber>;
        trafficLimitStrategy: z.ZodOptional<z.ZodEffects<z.ZodDefault<z.ZodNativeEnum<{
            readonly NO_RESET: "NO_RESET";
            readonly DAY: "DAY";
            readonly WEEK: "WEEK";
            readonly MONTH: "MONTH";
            readonly MONTH_ROLLING: "MONTH_ROLLING";
        }>>, "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING", "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING" | undefined>>;
        expireAt: z.ZodOptional<z.ZodEffects<z.ZodEffects<z.ZodString, Date, string>, Date, string>>;
        description: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        tag: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        telegramId: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
        email: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        hwidDeviceLimit: z.ZodOptional<z.ZodNullable<z.ZodNumber>>;
        activeInternalSquads: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
        externalSquadUuid: z.ZodOptional<z.ZodNullable<z.ZodString>>;
    }, "strip", z.ZodTypeAny, {
        status?: "DISABLED" | "ACTIVE" | undefined;
        uuid?: string | undefined;
        expireAt?: Date | undefined;
        description?: string | null | undefined;
        username?: string | undefined;
        tag?: string | null | undefined;
        trafficLimitBytes?: number | undefined;
        trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING" | undefined;
        telegramId?: number | null | undefined;
        email?: string | null | undefined;
        hwidDeviceLimit?: number | null | undefined;
        externalSquadUuid?: string | null | undefined;
        activeInternalSquads?: string[] | undefined;
    }, {
        status?: "DISABLED" | "ACTIVE" | undefined;
        uuid?: string | undefined;
        expireAt?: string | undefined;
        description?: string | null | undefined;
        username?: string | undefined;
        tag?: string | null | undefined;
        trafficLimitBytes?: number | undefined;
        trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING" | undefined;
        telegramId?: number | null | undefined;
        email?: string | null | undefined;
        hwidDeviceLimit?: number | null | undefined;
        externalSquadUuid?: string | null | undefined;
        activeInternalSquads?: string[] | undefined;
    }>, {
        status?: "DISABLED" | "ACTIVE" | undefined;
        uuid?: string | undefined;
        expireAt?: Date | undefined;
        description?: string | null | undefined;
        username?: string | undefined;
        tag?: string | null | undefined;
        trafficLimitBytes?: number | undefined;
        trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING" | undefined;
        telegramId?: number | null | undefined;
        email?: string | null | undefined;
        hwidDeviceLimit?: number | null | undefined;
        externalSquadUuid?: string | null | undefined;
        activeInternalSquads?: string[] | undefined;
    }, {
        status?: "DISABLED" | "ACTIVE" | undefined;
        uuid?: string | undefined;
        expireAt?: string | undefined;
        description?: string | null | undefined;
        username?: string | undefined;
        tag?: string | null | undefined;
        trafficLimitBytes?: number | undefined;
        trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING" | undefined;
        telegramId?: number | null | undefined;
        email?: string | null | undefined;
        hwidDeviceLimit?: number | null | undefined;
        externalSquadUuid?: string | null | undefined;
        activeInternalSquads?: string[] | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodString;
            id: z.ZodNumber;
            shortUuid: z.ZodString;
            username: z.ZodString;
            status: z.ZodDefault<z.ZodNativeEnum<{
                readonly ACTIVE: "ACTIVE";
                readonly DISABLED: "DISABLED";
                readonly LIMITED: "LIMITED";
                readonly EXPIRED: "EXPIRED";
            }>>;
            trafficLimitBytes: z.ZodDefault<z.ZodNumber>;
            trafficLimitStrategy: z.ZodDefault<z.ZodNativeEnum<{
                readonly NO_RESET: "NO_RESET";
                readonly DAY: "DAY";
                readonly WEEK: "WEEK";
                readonly MONTH: "MONTH";
                readonly MONTH_ROLLING: "MONTH_ROLLING";
            }>>;
            expireAt: z.ZodEffects<z.ZodString, Date, string>;
            telegramId: z.ZodNullable<z.ZodNumber>;
            email: z.ZodNullable<z.ZodString>;
            description: z.ZodNullable<z.ZodString>;
            tag: z.ZodNullable<z.ZodString>;
            hwidDeviceLimit: z.ZodNullable<z.ZodNumber>;
            externalSquadUuid: z.ZodNullable<z.ZodString>;
            trojanPassword: z.ZodString;
            vlessUuid: z.ZodString;
            ssPassword: z.ZodString;
            lastTriggeredThreshold: z.ZodDefault<z.ZodNumber>;
            subRevokedAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
            lastTrafficResetAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
        } & {
            subscriptionUrl: z.ZodString;
            activeInternalSquads: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                name: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                name: string;
            }, {
                uuid: string;
                name: string;
            }>, "many">;
            userTraffic: z.ZodObject<{
                usedTrafficBytes: z.ZodNumber;
                lifetimeUsedTrafficBytes: z.ZodNumber;
                onlineAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
                firstConnectedAt: z.ZodNullable<z.ZodEffects<z.ZodString, Date, string>>;
                lastConnectedNodeUuid: z.ZodNullable<z.ZodString>;
            }, "strip", z.ZodTypeAny, {
                usedTrafficBytes: number;
                lifetimeUsedTrafficBytes: number;
                onlineAt: Date | null;
                firstConnectedAt: Date | null;
                lastConnectedNodeUuid: string | null;
            }, {
                usedTrafficBytes: number;
                lifetimeUsedTrafficBytes: number;
                onlineAt: string | null;
                firstConnectedAt: string | null;
                lastConnectedNodeUuid: string | null;
            }>;
        }, "strip", z.ZodTypeAny, {
            status: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
            uuid: string;
            expireAt: Date;
            createdAt: Date;
            updatedAt: Date;
            description: string | null;
            username: string;
            tag: string | null;
            id: number;
            shortUuid: string;
            trafficLimitBytes: number;
            trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING";
            telegramId: number | null;
            email: string | null;
            hwidDeviceLimit: number | null;
            externalSquadUuid: string | null;
            trojanPassword: string;
            vlessUuid: string;
            ssPassword: string;
            lastTriggeredThreshold: number;
            subRevokedAt: Date | null;
            lastTrafficResetAt: Date | null;
            subscriptionUrl: string;
            activeInternalSquads: {
                uuid: string;
                name: string;
            }[];
            userTraffic: {
                usedTrafficBytes: number;
                lifetimeUsedTrafficBytes: number;
                onlineAt: Date | null;
                firstConnectedAt: Date | null;
                lastConnectedNodeUuid: string | null;
            };
        }, {
            uuid: string;
            expireAt: string;
            createdAt: string;
            updatedAt: string;
            description: string | null;
            username: string;
            tag: string | null;
            id: number;
            shortUuid: string;
            telegramId: number | null;
            email: string | null;
            hwidDeviceLimit: number | null;
            externalSquadUuid: string | null;
            trojanPassword: string;
            vlessUuid: string;
            ssPassword: string;
            subRevokedAt: string | null;
            lastTrafficResetAt: string | null;
            subscriptionUrl: string;
            activeInternalSquads: {
                uuid: string;
                name: string;
            }[];
            userTraffic: {
                usedTrafficBytes: number;
                lifetimeUsedTrafficBytes: number;
                onlineAt: string | null;
                firstConnectedAt: string | null;
                lastConnectedNodeUuid: string | null;
            };
            status?: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE" | undefined;
            trafficLimitBytes?: number | undefined;
            trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING" | undefined;
            lastTriggeredThreshold?: number | undefined;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            status: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE";
            uuid: string;
            expireAt: Date;
            createdAt: Date;
            updatedAt: Date;
            description: string | null;
            username: string;
            tag: string | null;
            id: number;
            shortUuid: string;
            trafficLimitBytes: number;
            trafficLimitStrategy: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING";
            telegramId: number | null;
            email: string | null;
            hwidDeviceLimit: number | null;
            externalSquadUuid: string | null;
            trojanPassword: string;
            vlessUuid: string;
            ssPassword: string;
            lastTriggeredThreshold: number;
            subRevokedAt: Date | null;
            lastTrafficResetAt: Date | null;
            subscriptionUrl: string;
            activeInternalSquads: {
                uuid: string;
                name: string;
            }[];
            userTraffic: {
                usedTrafficBytes: number;
                lifetimeUsedTrafficBytes: number;
                onlineAt: Date | null;
                firstConnectedAt: Date | null;
                lastConnectedNodeUuid: string | null;
            };
        };
    }, {
        response: {
            uuid: string;
            expireAt: string;
            createdAt: string;
            updatedAt: string;
            description: string | null;
            username: string;
            tag: string | null;
            id: number;
            shortUuid: string;
            telegramId: number | null;
            email: string | null;
            hwidDeviceLimit: number | null;
            externalSquadUuid: string | null;
            trojanPassword: string;
            vlessUuid: string;
            ssPassword: string;
            subRevokedAt: string | null;
            lastTrafficResetAt: string | null;
            subscriptionUrl: string;
            activeInternalSquads: {
                uuid: string;
                name: string;
            }[];
            userTraffic: {
                usedTrafficBytes: number;
                lifetimeUsedTrafficBytes: number;
                onlineAt: string | null;
                firstConnectedAt: string | null;
                lastConnectedNodeUuid: string | null;
            };
            status?: "DISABLED" | "LIMITED" | "EXPIRED" | "ACTIVE" | undefined;
            trafficLimitBytes?: number | undefined;
            trafficLimitStrategy?: "MONTH" | "NO_RESET" | "DAY" | "WEEK" | "MONTH_ROLLING" | undefined;
            lastTriggeredThreshold?: number | undefined;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-user.command.d.ts.map