import { z } from 'zod';
export declare namespace CloneNodePluginCommand {
    const url: "/api/node-plugins/actions/clone";
    const TSQ_url: "/api/node-plugins/actions/clone";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        cloneFromUuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        cloneFromUuid: string;
    }, {
        cloneFromUuid: string;
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
//# sourceMappingURL=clone-node-plugin.command.d.ts.map