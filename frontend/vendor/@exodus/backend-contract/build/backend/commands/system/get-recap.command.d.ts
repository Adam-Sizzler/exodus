import { z } from 'zod';
export declare namespace GetRecapCommand {
    const url: "/api/system/stats/recap";
    const TSQ_url: "/api/system/stats/recap";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            thisMonth: z.ZodObject<{
                users: z.ZodNumber;
                traffic: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                users: number;
                traffic: string;
            }, {
                users: number;
                traffic: string;
            }>;
            total: z.ZodObject<{
                users: z.ZodNumber;
                nodes: z.ZodNumber;
                traffic: z.ZodString;
                nodesRam: z.ZodString;
                nodesCpuCores: z.ZodNumber;
                distinctCountries: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                nodes: number;
                users: number;
                traffic: string;
                nodesRam: string;
                nodesCpuCores: number;
                distinctCountries: number;
            }, {
                nodes: number;
                users: number;
                traffic: string;
                nodesRam: string;
                nodesCpuCores: number;
                distinctCountries: number;
            }>;
            version: z.ZodString;
            initDate: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            total: {
                nodes: number;
                users: number;
                traffic: string;
                nodesRam: string;
                nodesCpuCores: number;
                distinctCountries: number;
            };
            version: string;
            thisMonth: {
                users: number;
                traffic: string;
            };
            initDate: Date;
        }, {
            total: {
                nodes: number;
                users: number;
                traffic: string;
                nodesRam: string;
                nodesCpuCores: number;
                distinctCountries: number;
            };
            version: string;
            thisMonth: {
                users: number;
                traffic: string;
            };
            initDate: string;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: {
                nodes: number;
                users: number;
                traffic: string;
                nodesRam: string;
                nodesCpuCores: number;
                distinctCountries: number;
            };
            version: string;
            thisMonth: {
                users: number;
                traffic: string;
            };
            initDate: Date;
        };
    }, {
        response: {
            total: {
                nodes: number;
                users: number;
                traffic: string;
                nodesRam: string;
                nodesCpuCores: number;
                distinctCountries: number;
            };
            version: string;
            thisMonth: {
                users: number;
                traffic: string;
            };
            initDate: string;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-recap.command.d.ts.map