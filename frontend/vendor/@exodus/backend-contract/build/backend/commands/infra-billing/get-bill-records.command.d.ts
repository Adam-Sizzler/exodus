import { z } from 'zod';
export declare namespace GetInfraBillingHistoryRecordsCommand {
    const url: "/api/infra-billing/history";
    const TSQ_url: "/api/infra-billing/history";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodNumber>;
        size: z.ZodDefault<z.ZodNumber>;
    }, "strip", z.ZodTypeAny, {
        start: number;
        size: number;
    }, {
        start?: number | undefined;
        size?: number | undefined;
    }>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            records: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                providerUuid: z.ZodString;
                amount: z.ZodNumber;
                billedAt: z.ZodEffects<z.ZodString, Date, string>;
                provider: z.ZodObject<Omit<{
                    uuid: z.ZodString;
                    name: z.ZodString;
                    faviconLink: z.ZodNullable<z.ZodString>;
                    loginUrl: z.ZodNullable<z.ZodString>;
                    createdAt: z.ZodEffects<z.ZodString, Date, string>;
                    updatedAt: z.ZodEffects<z.ZodString, Date, string>;
                }, "createdAt" | "updatedAt" | "loginUrl">, "strip", z.ZodTypeAny, {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                }, {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                }>;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                };
                providerUuid: string;
                amount: number;
                billedAt: Date;
            }, {
                uuid: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                };
                providerUuid: string;
                amount: number;
                billedAt: string;
            }>, "many">;
            total: z.ZodNumber;
        }, "strip", z.ZodTypeAny, {
            total: number;
            records: {
                uuid: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                };
                providerUuid: string;
                amount: number;
                billedAt: Date;
            }[];
        }, {
            total: number;
            records: {
                uuid: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                };
                providerUuid: string;
                amount: number;
                billedAt: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            total: number;
            records: {
                uuid: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                };
                providerUuid: string;
                amount: number;
                billedAt: Date;
            }[];
        };
    }, {
        response: {
            total: number;
            records: {
                uuid: string;
                provider: {
                    uuid: string;
                    name: string;
                    faviconLink: string | null;
                };
                providerUuid: string;
                amount: number;
                billedAt: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-bill-records.command.d.ts.map