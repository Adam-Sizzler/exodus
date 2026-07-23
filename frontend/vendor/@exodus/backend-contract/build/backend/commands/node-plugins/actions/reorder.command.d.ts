import { z } from 'zod';
export declare namespace ReorderNodePluginCommand {
    const url: "/api/node-plugins/actions/reorder";
    const TSQ_url: "/api/node-plugins/actions/reorder";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        items: z.ZodArray<z.ZodObject<Pick<{
            uuid: z.ZodString;
            viewPosition: z.ZodNumber;
            name: z.ZodString;
            pluginConfig: z.ZodNullable<z.ZodUnknown>;
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
            nodePlugins: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                viewPosition: z.ZodNumber;
                name: z.ZodString;
                pluginConfig: z.ZodNullable<z.ZodUnknown>;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                name: string;
                viewPosition: number;
                pluginConfig?: unknown;
            }, {
                uuid: string;
                name: string;
                viewPosition: number;
                pluginConfig?: unknown;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            total: number;
            nodePlugins: {
                uuid: string;
                name: string;
                viewPosition: number;
                pluginConfig?: unknown;
            }[];
        }, {
            total: number;
            nodePlugins: {
                uuid: string;
                name: string;
                viewPosition: number;
                pluginConfig?: unknown;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: number;
            nodePlugins: {
                uuid: string;
                name: string;
                viewPosition: number;
                pluginConfig?: unknown;
            }[];
        };
    }, {
        response: {
            total: number;
            nodePlugins: {
                uuid: string;
                name: string;
                viewPosition: number;
                pluginConfig?: unknown;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=reorder.command.d.ts.map