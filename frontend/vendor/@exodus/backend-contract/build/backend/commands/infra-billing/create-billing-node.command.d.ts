import { z } from 'zod';
export declare namespace CreateInfraBillingNodeCommand {
    const url: "/api/infra-billing/nodes";
    const TSQ_url: "/api/infra-billing/nodes";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        providerUuid: z.ZodUUID;
        nodeUuid: z.ZodNullable<z.ZodUUID>;
        name: z.ZodNullable<z.ZodString>;
        nextBillingAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            totalBillingNodes: z.ZodNumber;
            billingNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                nodeUuid: z.ZodNullable<z.ZodUUID>;
                name: z.ZodNullable<z.ZodString>;
                providerUuid: z.ZodUUID;
                provider: z.ZodObject<{
                    uuid: z.ZodUUID;
                    name: z.ZodString;
                    faviconLink: z.ZodNullable<z.ZodString>;
                    loginUrl: z.ZodNullable<z.ZodString>;
                }, z.core.$strip>;
                node: z.ZodNullable<z.ZodObject<{
                    uuid: z.ZodUUID;
                    name: z.ZodString;
                    countryCode: z.ZodString;
                }, z.core.$strip>>;
                nextBillingAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            }, z.core.$strip>>;
            availableBillingNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
                countryCode: z.ZodString;
            }, z.core.$strip>>;
            totalAvailableBillingNodes: z.ZodNumber;
            stats: z.ZodObject<{
                upcomingNodesCount: z.ZodNumber;
                currentMonthPayments: z.ZodNumber;
                totalSpent: z.ZodNumber;
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=create-billing-node.command.d.ts.map