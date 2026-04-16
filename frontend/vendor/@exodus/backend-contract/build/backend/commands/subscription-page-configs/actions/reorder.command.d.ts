import { z } from 'zod';
export declare namespace ReorderSubscriptionPageConfigsCommand {
    const url: "/api/subscription-page-configs/actions/reorder";
    const TSQ_url: "/api/subscription-page-configs/actions/reorder";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        items: z.ZodArray<z.ZodObject<Pick<{
            uuid: z.ZodString;
            viewPosition: z.ZodNumber;
            name: z.ZodString;
            config: z.ZodNullable<z.ZodUnknown>;
        }, "uuid" | "viewPosition">, "strip", z.ZodTypeAny, {
            uuid: string;
            viewPosition: number;
        }, {
            uuid: string;
            viewPosition: number;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        items: {
            uuid: string;
            viewPosition: number;
        }[];
    }, {
        items: {
            uuid: string;
            viewPosition: number;
        }[];
    }>;
    type Request = z.infer<typeof RequestSchema>;
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
//# sourceMappingURL=reorder.command.d.ts.map