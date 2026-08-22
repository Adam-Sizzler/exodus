import { z } from 'zod';
export declare namespace GetSharedListsCommand {
    const url: "/api/node-plugins/shared-lists";
    const TSQ_url: "/api/node-plugins/shared-lists";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            sharedLists: z.ZodArray<z.ZodObject<{
                name: z.ZodString;
                type: z.ZodString;
                itemsCount: z.ZodNumber;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-shared-lists.command.d.ts.map