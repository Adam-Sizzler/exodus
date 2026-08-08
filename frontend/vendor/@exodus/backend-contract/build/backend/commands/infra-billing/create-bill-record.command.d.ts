import { z } from 'zod';
export declare namespace CreateInfraBillingRecordCommand {
    const url: "/api/infra-billing/history";
    const TSQ_url: "/api/infra-billing/history";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        providerUuid: z.ZodUUID;
        amount: z.ZodNumber;
        billedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            records: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                providerUuid: z.ZodUUID;
                amount: z.ZodNumber;
                billedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                provider: z.ZodObject<{
                    uuid: z.ZodUUID;
                    name: z.ZodString;
                    faviconLink: z.ZodNullable<z.ZodString>;
                }, z.core.$strip>;
            }, z.core.$strip>>;
            total: z.ZodNumber;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create-bill-record.command.d.ts.map