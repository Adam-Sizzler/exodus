import z from 'zod';
export declare const CustomRemarksSchema: z.ZodObject<{
    expiredUsers: z.ZodArray<z.ZodString>;
    limitedUsers: z.ZodArray<z.ZodString>;
    disabledUsers: z.ZodArray<z.ZodString>;
    emptyHosts: z.ZodArray<z.ZodString>;
    HWIDMaxDevicesExceeded: z.ZodArray<z.ZodString>;
    HWIDNotSupported: z.ZodArray<z.ZodString>;
}, z.core.$strip>;
export type TCustomRemarks = z.infer<typeof CustomRemarksSchema>;
//# sourceMappingURL=custom-remarks.schema.d.ts.map