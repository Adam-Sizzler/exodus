import { z } from 'zod';
export declare namespace GetInfraProvidersCommand {
    const url: "/api/infra-billing/providers";
    const TSQ_url: "/api/infra-billing/providers";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            providers: z.ZodArray<z.ZodObject<{
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
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            providers: {
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
            }[];
            total: number;
        }, {
            providers: {
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
            }[];
            total: number;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            providers: {
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
            }[];
            total: number;
        };
    }, {
        response: {
            providers: {
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
            }[];
            total: number;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-infra-providers.command.d.ts.map