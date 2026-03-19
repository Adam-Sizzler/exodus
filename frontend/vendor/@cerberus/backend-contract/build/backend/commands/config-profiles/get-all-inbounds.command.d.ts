import { z } from 'zod';
export declare namespace GetAllInboundsCommand {
    const url: "/api/config-profiles/inbounds";
    const TSQ_url: "/api/config-profiles/inbounds";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            inbounds: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                profileUuid: z.ZodString;
                tag: z.ZodString;
                type: z.ZodString;
                network: z.ZodNullable<z.ZodString>;
                security: z.ZodNullable<z.ZodString>;
                port: z.ZodNullable<z.ZodNumber>;
                rawInbound: z.ZodNullable<z.ZodUnknown>;
            } & {
                activeSquads: z.ZodArray<z.ZodString, "many">;
            }, "strip", z.ZodTypeAny, {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                activeSquads: string[];
                rawInbound?: unknown;
            }, {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                activeSquads: string[];
                rawInbound?: unknown;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            inbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                activeSquads: string[];
                rawInbound?: unknown;
            }[];
            total: number;
        }, {
            inbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                activeSquads: string[];
                rawInbound?: unknown;
            }[];
            total: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            inbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                activeSquads: string[];
                rawInbound?: unknown;
            }[];
            total: number;
        };
    }, {
        response: {
            inbounds: {
                type: string;
                uuid: string;
                profileUuid: string;
                tag: string;
                network: string | null;
                security: string | null;
                port: number | null;
                activeSquads: string[];
                rawInbound?: unknown;
            }[];
            total: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-all-inbounds.command.d.ts.map