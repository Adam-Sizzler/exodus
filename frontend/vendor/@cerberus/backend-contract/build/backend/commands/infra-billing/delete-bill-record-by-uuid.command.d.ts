import { z } from 'zod';
export declare namespace DeleteInfraBillingHistoryRecordCommand {
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
//# sourceMappingURL=delete-bill-record-by-uuid.command.d.ts.map