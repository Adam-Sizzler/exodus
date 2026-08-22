import z from 'zod';
export declare const SharedListNameSchema: z.ZodString;
export declare const SharedListsSchema: z.ZodObject<{
    name: z.ZodString;
    config: z.ZodRecord<z.ZodString, z.ZodUnknown>;
}, z.core.$strip>;
export declare const SharedListPreviewSchema: z.ZodObject<{
    name: z.ZodString;
    type: z.ZodString;
    itemsCount: z.ZodNumber;
}, z.core.$strip>;
//# sourceMappingURL=shared-list.schema.d.ts.map