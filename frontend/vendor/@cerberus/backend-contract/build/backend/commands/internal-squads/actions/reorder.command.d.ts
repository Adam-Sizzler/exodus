import { z } from 'zod';
export declare namespace ReorderInternalSquadCommand {
    const url: "/api/internal-squads/actions/reorder";
    const TSQ_url: "/api/internal-squads/actions/reorder";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        items: z.ZodArray<z.ZodObject<Pick<{
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
        }, "uuid" | "viewPosition">, "strip", z.ZodTypeAny, {
            uuid: string;
            viewPosition: number;
        }, {
            uuid: string;
            viewPosition: number;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        items: {
            uuid: string;
            viewPosition: number;
        }[];
    }, {
        items: {
            uuid: string;
            viewPosition: number;
        }[];
    }>;
    type Request = z.infer<typeof RequestSchema>;
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
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            total: number;
            internalSquads: {
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
            }[];
        }, {
            total: number;
            internalSquads: {
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
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: number;
            internalSquads: {
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
            }[];
        };
    }, {
        response: {
            total: number;
            internalSquads: {
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
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=reorder.command.d.ts.map