import { z } from 'zod';
export declare namespace GetSubpageConfigByShortUuidCommand {
    const url: (shortUuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        shortUuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        shortUuid: string;
    }, {
        shortUuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const RequestBodySchema: z.ZodObject<{
        requestHeaders: z.ZodRecord<z.ZodString, z.ZodString>;
    }, "strip", z.ZodTypeAny, {
        requestHeaders: Record<string, string>;
    }, {
        requestHeaders: Record<string, string>;
    }>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            subpageConfigUuid: z.ZodNullable<z.ZodString>;
            webpageAllowed: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            subpageConfigUuid: string | null;
            webpageAllowed: boolean;
        }, {
            subpageConfigUuid: string | null;
            webpageAllowed: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            subpageConfigUuid: string | null;
            webpageAllowed: boolean;
        };
    }, {
        response: {
            subpageConfigUuid: string | null;
            webpageAllowed: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subpage-config-by-shortuuid.command.d.ts.map