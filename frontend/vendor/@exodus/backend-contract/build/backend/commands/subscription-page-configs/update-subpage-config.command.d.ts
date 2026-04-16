import { z } from 'zod';
export declare namespace UpdateSubscriptionPageConfigCommand {
    const url: "/api/subscription-page-configs/";
    const TSQ_url: "/api/subscription-page-configs/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
        name: z.ZodOptional<z.ZodString>;
        config: z.ZodOptional<z.ZodUnknown>;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
        name?: string | undefined;
        config?: unknown;
    }, {
        uuid: string;
        name?: string | undefined;
        config?: unknown;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodString;
            viewPosition: z.ZodNumber;
            name: z.ZodString;
        } & {
            config: z.ZodUnknown;
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
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            uuid: string;
            name: string;
            viewPosition: number;
            config?: unknown;
        };
    }, {
        response: {
            uuid: string;
            name: string;
            viewPosition: number;
            config?: unknown;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-subpage-config.command.d.ts.map