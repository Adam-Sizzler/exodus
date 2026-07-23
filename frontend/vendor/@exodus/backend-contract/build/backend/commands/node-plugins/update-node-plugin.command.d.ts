import { z } from 'zod';
export declare namespace UpdateNodePluginCommand {
    const url: "/api/node-plugins/";
    const TSQ_url: "/api/node-plugins/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
        name: z.ZodOptional<z.ZodString>;
        pluginConfig: z.ZodOptional<z.ZodUnknown>;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
        name?: string | undefined;
        pluginConfig?: unknown;
    }, {
        uuid: string;
        name?: string | undefined;
        pluginConfig?: unknown;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodString;
            viewPosition: z.ZodNumber;
            name: z.ZodString;
        } & {
            pluginConfig: z.ZodUnknown;
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
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            uuid: string;
            name: string;
            viewPosition: number;
            pluginConfig?: unknown;
        };
    }, {
        response: {
            uuid: string;
            name: string;
            viewPosition: number;
            pluginConfig?: unknown;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-node-plugin.command.d.ts.map