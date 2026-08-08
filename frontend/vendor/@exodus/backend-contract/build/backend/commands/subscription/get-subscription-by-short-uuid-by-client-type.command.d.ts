import { z } from 'zod';
export declare namespace GetSubscriptionByShortUuidByClientTypeCommand {
    const url: (shortUuid: string) => string;
    const TSQ_url: string;
    const RequestParamSchema: z.ZodObject<{
        shortUuid: z.ZodString;
        clientType: z.ZodEnum<{
            readonly STASH: "stash";
            readonly SINGBOX: "singbox";
            readonly MIHOMO: "mihomo";
            readonly XRAY_JSON: "json";
            readonly V2RAY_JSON: "v2ray-json";
            readonly CLASH: "clash";
        }>;
    }, z.core.$strip>;
    type RequestParam = z.infer<typeof RequestParamSchema>;
}
//# sourceMappingURL=get-subscription-by-short-uuid-by-client-type.command.d.ts.map