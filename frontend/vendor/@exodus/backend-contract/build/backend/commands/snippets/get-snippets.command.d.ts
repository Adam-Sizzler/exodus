import { z } from 'zod';
export declare namespace GetSnippetsCommand {
    const url: "/api/snippets/";
    const TSQ_url: "/api/snippets/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            snippets: z.ZodArray<z.ZodObject<{
                name: z.ZodString;
                snippet: z.ZodUnknown;
            }, "strip", z.ZodTypeAny, {
                name: string;
                snippet?: unknown;
            }, {
                name: string;
                snippet?: unknown;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            snippets: {
                name: string;
                snippet?: unknown;
            }[];
            total: number;
        }, {
            snippets: {
                name: string;
                snippet?: unknown;
            }[];
            total: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            snippets: {
                name: string;
                snippet?: unknown;
            }[];
            total: number;
        };
    }, {
        response: {
            snippets: {
                name: string;
                snippet?: unknown;
            }[];
            total: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-snippets.command.d.ts.map