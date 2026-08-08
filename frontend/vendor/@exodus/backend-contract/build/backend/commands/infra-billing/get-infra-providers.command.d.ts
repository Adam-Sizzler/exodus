import { z } from 'zod';
export declare namespace GetInfraProvidersCommand {
    const url: "/api/infra-billing/providers";
    const TSQ_url: "/api/infra-billing/providers";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            providers: z.ZodArray<z.ZodObject<{
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
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-infra-providers.command.d.ts.map