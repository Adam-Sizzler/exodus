import { z } from 'zod';
export declare namespace ReorderConfigProfileCommand {
    const url: "/api/config-profiles/actions/reorder";
    const TSQ_url: "/api/config-profiles/actions/reorder";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        items: z.ZodArray<z.ZodObject<Pick<{
            uuid: z.ZodString;
            viewPosition: z.ZodNumber;
            name: z.ZodString;
            config: z.ZodUnknown;
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
            nodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                name: z.ZodString;
                countryCode: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                countryCode: string;
                name: string;
            }, {
                uuid: string;
                countryCode: string;
                name: string;
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
            configProfiles: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                viewPosition: z.ZodNumber;
                name: z.ZodString;
                config: z.ZodUnknown;
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
                nodes: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodString;
                    name: z.ZodString;
                    countryCode: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    uuid: string;
                    countryCode: string;
                    name: string;
                }, {
                    uuid: string;
                    countryCode: string;
                    name: string;
                }>, "many">;
                createdAt: z.ZodEffects<z.ZodString, Date, string>;
                updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            }, "strip", z.ZodTypeAny, {
                nodes: {
                    uuid: string;
                    countryCode: string;
                    name: string;
                }[];
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
                config?: unknown;
            }, {
                nodes: {
                    uuid: string;
                    countryCode: string;
                    name: string;
                }[];
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
                config?: unknown;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            total: number;
            configProfiles: {
                nodes: {
                    uuid: string;
                    countryCode: string;
                    name: string;
                }[];
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
                config?: unknown;
            }[];
        }, {
            total: number;
            configProfiles: {
                nodes: {
                    uuid: string;
                    countryCode: string;
                    name: string;
                }[];
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
                config?: unknown;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: number;
            configProfiles: {
                nodes: {
                    uuid: string;
                    countryCode: string;
                    name: string;
                }[];
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
                config?: unknown;
            }[];
        };
    }, {
        response: {
            total: number;
            configProfiles: {
                nodes: {
                    uuid: string;
                    countryCode: string;
                    name: string;
                }[];
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
                config?: unknown;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=reorder.command.d.ts.map