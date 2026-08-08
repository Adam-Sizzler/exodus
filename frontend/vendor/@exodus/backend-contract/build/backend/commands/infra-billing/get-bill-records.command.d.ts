import { z } from 'zod';
export declare namespace GetInfraBillingRecordsCommand {
    const url: "/api/infra-billing/history";
    const TSQ_url: "/api/infra-billing/history";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
        size: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
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
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-bill-records.command.d.ts.map