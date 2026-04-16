import { z } from 'zod';
export declare const InfraProviderSchema: z.ZodObject<{
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
export declare const PartialInfraProviderSchema: z.ZodObject<{
    uuid: z.ZodString;
    name: z.ZodString;
    faviconLink: z.ZodNullable<z.ZodString>;
    loginUrl: z.ZodNullable<z.ZodString>;
    createdAt: z.ZodEffects<z.ZodString, Date, string>;
    updatedAt: z.ZodEffects<z.ZodString, Date, string>;
}, "strip", z.ZodTypeAny, {
    uuid: string;
    createdAt: Date;
    updatedAt: Date;
    name: string;
    faviconLink: string | null;
    loginUrl: string | null;
}, {
    uuid: string;
    createdAt: string;
    updatedAt: string;
    name: string;
    faviconLink: string | null;
    loginUrl: string | null;
}>;
//# sourceMappingURL=infra-provider.schema.d.ts.map