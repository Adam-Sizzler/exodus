import z from 'zod';
export declare const CustomRemarksSchema: z.ZodObject<{
    expiredUsers: z.ZodArray<z.ZodString, "many">;
    limitedUsers: z.ZodArray<z.ZodString, "many">;
    disabledUsers: z.ZodArray<z.ZodString, "many">;
    emptyHosts: z.ZodArray<z.ZodString, "many">;
    HWIDMaxDevicesExceeded: z.ZodArray<z.ZodString, "many">;
    HWIDNotSupported: z.ZodArray<z.ZodString, "many">;
}, "strip", z.ZodTypeAny, {
    expiredUsers: string[];
    limitedUsers: string[];
    disabledUsers: string[];
    emptyHosts: string[];
    HWIDMaxDevicesExceeded: string[];
    HWIDNotSupported: string[];
}, {
    expiredUsers: string[];
    limitedUsers: string[];
    disabledUsers: string[];
    emptyHosts: string[];
    HWIDMaxDevicesExceeded: string[];
    HWIDNotSupported: string[];
}>;
export type TCustomRemarks = z.infer<typeof CustomRemarksSchema>;
//# sourceMappingURL=custom-remarks.schema.d.ts.map