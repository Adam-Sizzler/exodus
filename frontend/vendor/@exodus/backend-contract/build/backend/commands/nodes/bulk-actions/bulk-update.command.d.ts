import { z } from 'zod';
export declare namespace BulkNodesUpdateCommand {
    const url: "/api/nodes/bulk-actions/update";
    const TSQ_url: "/api/nodes/bulk-actions/update";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuids: z.ZodArray<z.ZodString, "many">;
        fields: z.ZodObject<{
            countryCode: z.ZodOptional<z.ZodString>;
            consumptionMultiplier: z.ZodOptional<z.ZodEffects<z.ZodNumber, number, number>>;
            providerUuid: z.ZodOptional<z.ZodNullable<z.ZodString>>;
            tags: z.ZodOptional<z.ZodArray<z.ZodString, "many">>;
            activePluginUuid: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        }, "strip", z.ZodTypeAny, {
            tags?: string[] | undefined;
            countryCode?: string | undefined;
            consumptionMultiplier?: number | undefined;
            providerUuid?: string | null | undefined;
            activePluginUuid?: string | null | undefined;
        }, {
            tags?: string[] | undefined;
            countryCode?: string | undefined;
            consumptionMultiplier?: number | undefined;
            providerUuid?: string | null | undefined;
            activePluginUuid?: string | null | undefined;
        }>;
    }, "strip", z.ZodTypeAny, {
        uuids: string[];
        fields: {
            tags?: string[] | undefined;
            countryCode?: string | undefined;
            consumptionMultiplier?: number | undefined;
            providerUuid?: string | null | undefined;
            activePluginUuid?: string | null | undefined;
        };
    }, {
        uuids: string[];
        fields: {
            tags?: string[] | undefined;
            countryCode?: string | undefined;
            consumptionMultiplier?: number | undefined;
            providerUuid?: string | null | undefined;
            activePluginUuid?: string | null | undefined;
        };
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            eventSent: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            eventSent: boolean;
        }, {
            eventSent: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            eventSent: boolean;
        };
    }, {
        response: {
            eventSent: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=bulk-update.command.d.ts.map