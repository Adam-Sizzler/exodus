import { z } from 'zod';
export declare namespace CloneSubscriptionPageConfigCommand {
    const url: "/api/subscription-page-configs/actions/clone";
    const TSQ_url: "/api/subscription-page-configs/actions/clone";
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
//# sourceMappingURL=clone-subpage-config.command.d.ts.map