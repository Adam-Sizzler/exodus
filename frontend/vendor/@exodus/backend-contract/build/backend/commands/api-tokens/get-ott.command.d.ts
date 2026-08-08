import { z } from 'zod';
export declare namespace GetOttCommand {
    const url: "/api/tokens/ott";
    const TSQ_url: "/api/tokens/ott";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            ott: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-ott.command.d.ts.map