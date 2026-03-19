import { z } from 'zod';
export declare namespace GetSubscriptionPageConfigsCommand {
    const url: "/api/subscription-page-configs/";
    const TSQ_url: "/api/subscription-page-configs/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            configs: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                viewPosition: z.ZodNumber;
                name: z.ZodString;
                config: z.ZodNullable<z.ZodUnknown>;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                name: string;
                viewPosition: number;
                config?: unknown;
            }, {
                uuid: string;
                name: string;
                viewPosition: number;
                config?: unknown;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            total: number;
            configs: {
                uuid: string;
                name: string;
                viewPosition: number;
                config?: unknown;
            }[];
        }, {
            total: number;
            configs: {
                uuid: string;
                name: string;
                viewPosition: number;
                config?: unknown;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: number;
            configs: {
                uuid: string;
                name: string;
                viewPosition: number;
                config?: unknown;
            }[];
        };
    }, {
        response: {
            total: number;
            configs: {
                uuid: string;
                name: string;
                viewPosition: number;
                config?: unknown;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-subpage-configs.command.d.ts.map