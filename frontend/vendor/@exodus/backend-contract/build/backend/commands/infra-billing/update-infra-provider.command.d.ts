import { z } from 'zod';
export declare namespace UpdateInfraProviderCommand {
    const url: "/api/infra-billing/providers";
    const TSQ_url: "/api/infra-billing/providers";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
        name: z.ZodOptional<z.ZodString>;
        faviconLink: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        loginUrl: z.ZodOptional<z.ZodNullable<z.ZodString>>;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
        name?: string | undefined;
        faviconLink?: string | null | undefined;
        loginUrl?: string | null | undefined;
    }, {
        uuid: string;
        name?: string | undefined;
        faviconLink?: string | null | undefined;
        loginUrl?: string | null | undefined;
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
                name: z.ZodString;
                details: z.ZodNullable<z.ZodObject<{
                    nodeUuid: z.ZodString;
                    countryCode: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    nodeUuid: string;
                    countryCode: string;
                }, {
                    nodeUuid: string;
                    countryCode: string;
                }>>;
            }, "strip", z.ZodTypeAny, {
                name: string;
                details: {
                    nodeUuid: string;
                    countryCode: string;
                } | null;
            }, {
                name: string;
                details: {
                    nodeUuid: string;
                    countryCode: string;
                } | null;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            faviconLink: string | null;
            loginUrl: string | null;
            billingHistory: {
                totalAmount: number;
                totalBills: number;
            };
            billingNodes: {
                name: string;
                details: {
                    nodeUuid: string;
                    countryCode: string;
                } | null;
            }[];
        }, {
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            faviconLink: string | null;
            loginUrl: string | null;
            billingHistory: {
                totalAmount: number;
                totalBills: number;
            };
            billingNodes: {
                name: string;
                details: {
                    nodeUuid: string;
                    countryCode: string;
                } | null;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            uuid: string;
            name: string;
            createdAt: Date;
            updatedAt: Date;
            faviconLink: string | null;
            loginUrl: string | null;
            billingHistory: {
                totalAmount: number;
                totalBills: number;
            };
            billingNodes: {
                name: string;
                details: {
                    nodeUuid: string;
                    countryCode: string;
                } | null;
            }[];
        };
    }, {
        response: {
            uuid: string;
            name: string;
            createdAt: string;
            updatedAt: string;
            faviconLink: string | null;
            loginUrl: string | null;
            billingHistory: {
                totalAmount: number;
                totalBills: number;
            };
            billingNodes: {
                name: string;
                details: {
                    nodeUuid: string;
                    countryCode: string;
                } | null;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-infra-provider.command.d.ts.map