import { z } from 'zod';
export declare namespace GetLegacyStatsNodeUserUsageCommand {
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
            username: z.ZodString;
            nodeUuid: z.ZodString;
            total: z.ZodNumber;
            date: z.ZodEffects<z.ZodString, Date, string>;
        }, "strip", z.ZodTypeAny, {
            date: Date;
            username: string;
            userUuid: string;
            nodeUuid: string;
            total: number;
        }, {
            date: string;
            username: string;
            userUuid: string;
            nodeUuid: string;
            total: number;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        response: {
            date: Date;
            username: string;
            userUuid: string;
            nodeUuid: string;
            total: number;
        }[];
    }, {
        response: {
            date: string;
            username: string;
            userUuid: string;
            nodeUuid: string;
            total: number;
        }[];
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-legacy-stats-node-user-usage.command.d.ts.map