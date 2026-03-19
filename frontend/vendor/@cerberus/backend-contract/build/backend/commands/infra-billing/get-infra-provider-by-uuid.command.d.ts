import { z } from 'zod';
export declare namespace GetInfraProviderByUuidCommand {
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
            uuid: z.ZodString;
            name: z.ZodString;
            faviconLink: z.ZodNullable<z.ZodString>;
            loginUrl: z.ZodNullable<z.ZodString>;
            createdAt: z.ZodEffects<z.ZodString, Date, string>;
            updatedAt: z.ZodEffects<z.ZodString, Date, string>;
            billingHistory: z.ZodObject<{
                totalAmount: z.ZodNumber;
                totalBills: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                totalAmount: number;
                totalBills: number;
            }, {
                totalAmount: number;
                totalBills: number;
            }>;
            billingNodes: z.ZodArray<z.ZodObject<{
                nodeUuid: z.ZodString;
                name: z.ZodString;
                countryCode: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                nodeUuid: string;
                countryCode: string;
                name: string;
            }, {
                nodeUuid: string;
                countryCode: string;
                name: string;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
            billingHistory: {
                totalAmount: number;
                totalBills: number;
            };
            billingNodes: {
                nodeUuid: string;
                countryCode: string;
                name: string;
            }[];
        }, {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
            billingHistory: {
                totalAmount: number;
                totalBills: number;
            };
            billingNodes: {
                nodeUuid: string;
                countryCode: string;
                name: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            uuid: string;
            createdAt: Date;
            updatedAt: Date;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
            billingHistory: {
                totalAmount: number;
                totalBills: number;
            };
            billingNodes: {
                nodeUuid: string;
                countryCode: string;
                name: string;
            }[];
        };
    }, {
        response: {
            uuid: string;
            createdAt: string;
            updatedAt: string;
            name: string;
            faviconLink: string | null;
            loginUrl: string | null;
            billingHistory: {
                totalAmount: number;
                totalBills: number;
            };
            billingNodes: {
                nodeUuid: string;
                countryCode: string;
                name: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-infra-provider-by-uuid.command.d.ts.map