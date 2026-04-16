import { z } from 'zod';
export declare namespace GetExodusHealthCommand {
    const url: "/api/system/health";
    const TSQ_url: "/api/system/health";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            pm2Stats: z.ZodArray<z.ZodObject<{
                name: z.ZodString;
                memory: z.ZodString;
                cpu: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                name: string;
                memory: string;
                cpu: string;
            }, {
                name: string;
                memory: string;
                cpu: string;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            pm2Stats: {
                name: string;
                memory: string;
                cpu: string;
            }[];
        }, {
            pm2Stats: {
                name: string;
                memory: string;
                cpu: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            pm2Stats: {
                name: string;
                memory: string;
                cpu: string;
            }[];
        };
    }, {
        response: {
            pm2Stats: {
                name: string;
                memory: string;
                cpu: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
