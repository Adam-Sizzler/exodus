import { z } from 'zod';
export declare namespace GetNodeMetadataCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestParamsSchema: z.ZodObject<{
        uuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
    }, {
        uuid: string;
    }>;
    type RequestParams = z.infer<typeof RequestParamsSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            metadata: z.ZodObject<{}, "passthrough", z.ZodTypeAny, z.objectOutputType<{}, z.ZodTypeAny, "passthrough">, z.objectInputType<{}, z.ZodTypeAny, "passthrough">>;
        }, "strip", z.ZodTypeAny, {
            metadata: {} & {
                [k: string]: unknown;
            };
        }, {
            metadata: {} & {
                [k: string]: unknown;
            };
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            metadata: {} & {
                [k: string]: unknown;
            };
        };
    }, {
        response: {
            metadata: {} & {
                [k: string]: unknown;
            };
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-node-metadata.command.d.ts.map