import { z } from 'zod';
export declare namespace GetNodesMetricsCommand {
    const url: "/api/system/nodes/metrics";
    const TSQ_url: "/api/system/nodes/metrics";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            nodes: z.ZodArray<z.ZodObject<{
                nodeUuid: z.ZodString;
                nodeName: z.ZodString;
                countryEmoji: z.ZodString;
                providerName: z.ZodString;
                usersOnline: z.ZodNumber;
                inboundsStats: z.ZodArray<z.ZodObject<{
                    tag: z.ZodString;
                    upload: z.ZodString;
                    download: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    tag: string;
                    upload: string;
                    download: string;
                }, {
                    tag: string;
                    upload: string;
                    download: string;
                }>, "many">;
                outboundsStats: z.ZodArray<z.ZodObject<{
                    tag: z.ZodString;
                    upload: z.ZodString;
                    download: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    tag: string;
                    upload: string;
                    download: string;
                }, {
                    tag: string;
                    upload: string;
                    download: string;
                }>, "many">;
            }, "strip", z.ZodTypeAny, {
                nodeUuid: string;
                nodeName: string;
                usersOnline: number;
                providerName: string;
                countryEmoji: string;
                inboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
                outboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
            }, {
                nodeUuid: string;
                nodeName: string;
                usersOnline: number;
                providerName: string;
                countryEmoji: string;
                inboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
                outboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            nodes: {
                nodeUuid: string;
                nodeName: string;
                usersOnline: number;
                providerName: string;
                countryEmoji: string;
                inboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
                outboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
            }[];
        }, {
            nodes: {
                nodeUuid: string;
                nodeName: string;
                usersOnline: number;
                providerName: string;
                countryEmoji: string;
                inboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
                outboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            nodes: {
                nodeUuid: string;
                nodeName: string;
                usersOnline: number;
                providerName: string;
                countryEmoji: string;
                inboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
                outboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
            }[];
        };
    }, {
        response: {
            nodes: {
                nodeUuid: string;
                nodeName: string;
                usersOnline: number;
                providerName: string;
                countryEmoji: string;
                inboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
                outboundsStats: {
                    tag: string;
                    upload: string;
                    download: string;
                }[];
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-nodes-metrics.command.d.ts.map