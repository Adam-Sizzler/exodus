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
            }, z.core.$strip>;
            total: z.ZodObject<{
                users: z.ZodNumber;
                nodes: z.ZodNumber;
                traffic: z.ZodString;
                nodesRam: z.ZodString;
                nodesCpuCores: z.ZodNumber;
                distinctCountries: z.ZodNumber;
            }, z.core.$strip>;
            version: z.ZodString;
            initDate: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-recap.command.d.ts.map