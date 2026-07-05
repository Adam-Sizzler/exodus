import { z } from 'zod';
export declare namespace GetNodePluginsCommand {
    const url: "/api/node-plugins/";
    const TSQ_url: "/api/node-plugins/";
    const endpointDetails: import("../../constants").EndpointDetails;
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
//# sourceMappingURL=get-node-plugins.command.d.ts.map