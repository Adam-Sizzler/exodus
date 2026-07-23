import { z } from 'zod';
export declare namespace GetInternalSquadsCommand {
    const url: "/api/internal-squads/";
    const TSQ_url: "/api/internal-squads/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            internalSquads: z.ZodArray<z.ZodObject<{
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
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }, {
                    uuid: string;
                    type: string;
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
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                viewPosition: number;
                info: {
                    membersCount: number;
                    inboundsCount: number;
                };
            }, {
                inbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                viewPosition: number;
                info: {
                    membersCount: number;
                    inboundsCount: number;
                };
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            total: number;
            internalSquads: {
                inbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                viewPosition: number;
                info: {
                    membersCount: number;
                    inboundsCount: number;
                };
            }[];
        }, {
            total: number;
            internalSquads: {
                inbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                viewPosition: number;
                info: {
                    membersCount: number;
                    inboundsCount: number;
                };
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: number;
            internalSquads: {
                inbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
                uuid: string;
                name: string;
                createdAt: Date;
                updatedAt: Date;
                viewPosition: number;
                info: {
                    membersCount: number;
                    inboundsCount: number;
                };
            }[];
        };
    }, {
        response: {
            total: number;
            internalSquads: {
                inbounds: {
                    uuid: string;
                    type: string;
                    profileUuid: string;
                    tag: string;
                    network: string | null;
                    security: string | null;
                    port: number | null;
                    rawInbound?: unknown;
                }[];
                uuid: string;
                name: string;
                createdAt: string;
                updatedAt: string;
                viewPosition: number;
                info: {
                    membersCount: number;
                    inboundsCount: number;
                };
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-internal-squads.command.d.ts.map