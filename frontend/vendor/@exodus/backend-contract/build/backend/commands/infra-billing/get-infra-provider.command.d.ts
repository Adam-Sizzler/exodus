import { z } from 'zod';
export declare namespace GetInfraProviderCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        uuid: z.ZodUUID;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            name: z.ZodString;
            faviconLink: z.ZodNullable<z.ZodString>;
            loginUrl: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            billingHistory: z.ZodObject<{
                totalAmount: z.ZodNumber;
                totalBills: z.ZodNumber;
            }, z.core.$strip>;
            billingNodes: z.ZodArray<z.ZodObject<{
                name: z.ZodString;
                details: z.ZodNullable<z.ZodObject<{
                    nodeUuid: z.ZodUUID;
                    countryCode: z.ZodString;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-infra-provider.command.d.ts.map