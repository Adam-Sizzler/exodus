import { z } from 'zod';
export declare namespace GetMetadataCommand {
    const url: "/api/system/metadata";
    const TSQ_url: "/api/system/metadata";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            version: z.ZodString;
            build: z.ZodObject<{
                time: z.ZodString;
                number: z.ZodString;
            }, z.core.$strip>;
            git: z.ZodObject<{
                backend: z.ZodObject<{
                    commitSha: z.ZodString;
                    branch: z.ZodString;
                    commitUrl: z.ZodString;
                }, z.core.$strip>;
                frontend: z.ZodObject<{
                    commitSha: z.ZodString;
                    commitUrl: z.ZodString;
                }, z.core.$strip>;
            }, z.core.$strip>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-metadata.command.d.ts.map