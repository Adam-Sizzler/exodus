import { z } from 'zod';
export declare namespace DeleteHostCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
    }, {
        uuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            isDeleted: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            isDeleted: boolean;
        }, {
            isDeleted: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            isDeleted: boolean;
        };
    }, {
        response: {
            isDeleted: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=delete.command.d.ts.map