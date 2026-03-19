import { z } from 'zod';
export declare namespace GetLegacyStatsUserUsageCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
    }, {
        uuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodString;
        end: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        start: string;
        end: string;
    }, {
        start: string;
        end: string;
    }>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodArray<z.ZodObject<{
            userUuid: z.ZodString;
            nodeUuid: z.ZodString;
            nodeName: z.ZodString;
            countryCode: z.ZodString;
            total: z.ZodNumber;
            date: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            date: Date;
            userUuid: string;
            nodeUuid: string;
            total: number;
            nodeName: string;
            countryCode: string;
        }, {
            date: string;
            userUuid: string;
            nodeUuid: string;
            total: number;
            nodeName: string;
            countryCode: string;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        response: {
            date: Date;
            userUuid: string;
            nodeUuid: string;
            total: number;
            nodeName: string;
            countryCode: string;
        }[];
    }, {
        response: {
            date: string;
            userUuid: string;
            nodeUuid: string;
            total: number;
            nodeName: string;
            countryCode: string;
        }[];
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-legacy-user-usage.command.d.ts.map