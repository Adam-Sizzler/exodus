import { z } from 'zod';
export declare namespace GetSubscriptionByShortUuidByClientTypeCommand {
    const url: (shortUuid: string) => string;
    const TSQ_url: string;
    const RequestSchema: z.ZodObject<{
        shortUuid: z.ZodString;
        clientType: z.ZodNativeEnum<{
            readonly STASH: "stash";
            readonly SINGBOX: "singbox";
            readonly MIHOMO: "mihomo";
            readonly XRAY_JSON: "json";
            readonly V2RAY_JSON: "v2ray-json";
            readonly CLASH: "clash";
        }>;
    }, "strip", z.ZodTypeAny, {
        shortUuid: string;
        clientType: "stash" | "singbox" | "mihomo" | "json" | "v2ray-json" | "clash";
    }, {
        shortUuid: string;
        clientType: "stash" | "singbox" | "mihomo" | "json" | "v2ray-json" | "clash";
    }>;
    type Request = z.infer<typeof RequestSchema>;
}
//# sourceMappingURL=get-subscription-by-short-uuid-by-client-type.command.d.ts.map