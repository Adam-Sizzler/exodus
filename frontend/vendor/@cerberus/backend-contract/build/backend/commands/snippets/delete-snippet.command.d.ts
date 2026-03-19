import { z } from 'zod';
export declare namespace DeleteSnippetCommand {
    const url: "/api/snippets/";
    const TSQ_url: "/api/snippets/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        name: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        name: string;
    }, {
        name: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
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
//# sourceMappingURL=delete-snippet.command.d.ts.map