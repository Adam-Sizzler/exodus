import { z } from 'zod';
export declare namespace UpdateInfraProviderCommand {
    const url: "/api/infra-billing/providers";
    const TSQ_url: "/api/infra-billing/providers";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodOptional<z.ZodString>;
        faviconLink: z.ZodOptional<z.ZodNullable<z.ZodURL>>;
        loginUrl: z.ZodOptional<z.ZodNullable<z.ZodURL>>;
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
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-infra-provider.command.d.ts.map