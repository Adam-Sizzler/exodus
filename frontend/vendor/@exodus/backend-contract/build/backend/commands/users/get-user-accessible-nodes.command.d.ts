import { z } from 'zod';
export declare namespace GetUserAccessibleNodesCommand {
    const url: (userId: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        userId: z.ZodCoercedNumber<unknown>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            userId: z.ZodNumber;
            activeNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                nodeName: z.ZodString;
                countryCode: z.ZodString;
                configProfileUuid: z.ZodUUID;
                configProfileName: z.ZodString;
                activeSquads: z.ZodArray<z.ZodObject<{
                    squadName: z.ZodString;
                    activeInbounds: z.ZodArray<z.ZodString>;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-user-accessible-nodes.command.d.ts.map