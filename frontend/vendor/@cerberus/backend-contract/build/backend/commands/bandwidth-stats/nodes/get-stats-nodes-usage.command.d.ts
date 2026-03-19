import { z } from 'zod';
export declare namespace GetStatsNodesUsageCommand {
    const url: "/api/bandwidth-stats/nodes/";
    const TSQ_url: "/api/bandwidth-stats/nodes/";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodString;
        end: z.ZodString;
        topNodesLimit: z.ZodDefault<z.ZodNumber>;
    }, "strip", z.ZodTypeAny, {
        start: string;
        end: string;
        topNodesLimit: number;
    }, {
        start: string;
        end: string;
        topNodesLimit?: number | undefined;
    }>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            categories: z.ZodArray<z.ZodString, "many">;
            sparklineData: z.ZodArray<z.ZodNumber, "many">;
            topNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                color: z.ZodString;
                name: z.ZodString;
                countryCode: z.ZodString;
                total: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }, {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }>, "many">;
            series: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                name: z.ZodString;
                color: z.ZodString;
                countryCode: z.ZodString;
                total: z.ZodNumber;
                data: z.ZodArray<z.ZodNumber, "many">;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
                data: number[];
            }, {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
                data: number[];
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            categories: string[];
            sparklineData: number[];
            topNodes: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }[];
            series: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
                data: number[];
            }[];
        }, {
            categories: string[];
            sparklineData: number[];
            topNodes: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }[];
            series: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
                data: number[];
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            categories: string[];
            sparklineData: number[];
            topNodes: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }[];
            series: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
                data: number[];
            }[];
        };
    }, {
        response: {
            categories: string[];
            sparklineData: number[];
            topNodes: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }[];
            series: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
                data: number[];
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-stats-nodes-usage.command.d.ts.map