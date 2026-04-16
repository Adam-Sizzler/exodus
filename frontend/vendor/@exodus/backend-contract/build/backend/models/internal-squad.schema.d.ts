import { z } from 'zod';
export declare const InternalSquadSchema: z.ZodObject<{
    uuid: z.ZodString;
    viewPosition: z.ZodNumber;
    name: z.ZodString;
    info: z.ZodObject<{
        membersCount: z.ZodNumber;
        inboundsCount: z.ZodNumber;
    }, "strip", z.ZodTypeAny, {
        membersCount: number;
        inboundsCount: number;
    }, {
        membersCount: number;
        inboundsCount: number;
    }>;
    inbounds: z.ZodArray<z.ZodObject<{
        uuid: z.ZodString;
        profileUuid: z.ZodString;
        tag: z.ZodString;
        type: z.ZodString;
        network: z.ZodNullable<z.ZodString>;
        security: z.ZodNullable<z.ZodString>;
        port: z.ZodNullable<z.ZodNumber>;
        rawInbound: z.ZodNullable<z.ZodUnknown>;
    }, "strip", z.ZodTypeAny, {
        type: string;
        uuid: string;
        profileUuid: string;
        tag: string;
        network: string | null;
        security: string | null;
        port: number | null;
        rawInbound?: unknown;
    }, {
        type: string;
        uuid: string;
        profileUuid: string;
        tag: string;
        network: string | null;
        security: string | null;
        port: number | null;
        rawInbound?: unknown;
    }>, "many">;
    createdAt: z.ZodEffects<z.ZodString, Date, string>;
    updatedAt: z.ZodEffects<z.ZodString, Date, string>;
}, "strip", z.ZodTypeAny, {
    inbounds: {
        type: string;
        uuid: string;
        profileUuid: string;
        tag: string;
        network: string | null;
        security: string | null;
        port: number | null;
        rawInbound?: unknown;
    }[];
    uuid: string;
    createdAt: Date;
    updatedAt: Date;
    name: string;
    viewPosition: number;
    info: {
        membersCount: number;
        inboundsCount: number;
    };
}, {
    inbounds: {
        type: string;
        uuid: string;
        profileUuid: string;
        tag: string;
        network: string | null;
        security: string | null;
        port: number | null;
        rawInbound?: unknown;
    }[];
    uuid: string;
    createdAt: string;
    updatedAt: string;
    name: string;
    viewPosition: number;
    info: {
        membersCount: number;
        inboundsCount: number;
    };
}>;
//# sourceMappingURL=internal-squad.schema.d.ts.map