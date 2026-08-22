import { z } from 'zod';
export declare namespace DeleteSharedListCommand {
    const url: (name: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        name: z.ZodString;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
}
//# sourceMappingURL=delete-shared-list.command.d.ts.map