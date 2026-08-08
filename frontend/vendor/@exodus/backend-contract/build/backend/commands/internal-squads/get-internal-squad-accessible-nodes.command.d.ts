import { z } from 'zod';
export declare namespace GetInternalSquadAccessibleNodesCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        uuid: z.ZodUUID;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            squadUuid: z.ZodUUID;
            accessibleNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                nodeName: z.ZodString;
                countryCode: z.ZodString;
                configProfileUuid: z.ZodUUID;
                configProfileName: z.ZodString;
                activeInbounds: z.ZodArray<z.ZodString>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-internal-squad-accessible-nodes.command.d.ts.map