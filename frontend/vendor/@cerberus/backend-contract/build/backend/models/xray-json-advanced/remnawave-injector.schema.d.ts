import z from 'zod';
declare const HostSelectorSchema: z.ZodDiscriminatedUnion<"type", [z.ZodObject<{
    type: z.ZodLiteral<"uuids">;
    values: z.ZodArray<z.ZodString, "many">;
}, "strip", z.ZodTypeAny, {
    values: string[];
    type: "uuids";
}, {
    values: string[];
    type: "uuids";
}>, z.ZodObject<{
    type: z.ZodLiteral<"remarkRegex">;
    pattern: z.ZodString;
}, "strip", z.ZodTypeAny, {
    type: "remarkRegex";
    pattern: string;
}, {
    type: "remarkRegex";
    pattern: string;
}>, z.ZodObject<{
    type: z.ZodLiteral<"tagRegex">;
    pattern: z.ZodString;
}, "strip", z.ZodTypeAny, {
    type: "tagRegex";
    pattern: string;
}, {
    type: "tagRegex";
    pattern: string;
}>, z.ZodObject<{
    type: z.ZodLiteral<"sameTagAsRecipient">;
}, "strip", z.ZodTypeAny, {
    type: "sameTagAsRecipient";
}, {
    type: "sameTagAsRecipient";
}>]>;
declare const InjectHostsEntrySchema: z.ZodObject<{
    selector: z.ZodDiscriminatedUnion<"type", [z.ZodObject<{
        type: z.ZodLiteral<"uuids">;
        values: z.ZodArray<z.ZodString, "many">;
    }, "strip", z.ZodTypeAny, {
        values: string[];
        type: "uuids";
    }, {
        values: string[];
        type: "uuids";
    }>, z.ZodObject<{
        type: z.ZodLiteral<"remarkRegex">;
        pattern: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        type: "remarkRegex";
        pattern: string;
    }, {
        type: "remarkRegex";
        pattern: string;
    }>, z.ZodObject<{
        type: z.ZodLiteral<"tagRegex">;
        pattern: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        type: "tagRegex";
        pattern: string;
    }, {
        type: "tagRegex";
        pattern: string;
    }>, z.ZodObject<{
        type: z.ZodLiteral<"sameTagAsRecipient">;
    }, "strip", z.ZodTypeAny, {
        type: "sameTagAsRecipient";
    }, {
        type: "sameTagAsRecipient";
    }>]>;
    selectFrom: z.ZodOptional<z.ZodEnum<["ALL", "HIDDEN", "NOT_HIDDEN"]>>;
    tagPrefix: z.ZodOptional<z.ZodString>;
    useHostRemarkAsTag: z.ZodOptional<z.ZodBoolean>;
    useHostTagAsTag: z.ZodOptional<z.ZodBoolean>;
}, "strip", z.ZodTypeAny, {
    selector: {
        values: string[];
        type: "uuids";
    } | {
        type: "remarkRegex";
        pattern: string;
    } | {
        type: "tagRegex";
        pattern: string;
    } | {
        type: "sameTagAsRecipient";
    };
    selectFrom?: "ALL" | "HIDDEN" | "NOT_HIDDEN" | undefined;
    tagPrefix?: string | undefined;
    useHostRemarkAsTag?: boolean | undefined;
    useHostTagAsTag?: boolean | undefined;
}, {
    selector: {
        values: string[];
        type: "uuids";
    } | {
        type: "remarkRegex";
        pattern: string;
    } | {
        type: "tagRegex";
        pattern: string;
    } | {
        type: "sameTagAsRecipient";
    };
    selectFrom?: "ALL" | "HIDDEN" | "NOT_HIDDEN" | undefined;
    tagPrefix?: string | undefined;
    useHostRemarkAsTag?: boolean | undefined;
    useHostTagAsTag?: boolean | undefined;
}>;
export declare const CerberusInjectorSchema: z.ZodObject<{
    injectHosts: z.ZodOptional<z.ZodArray<z.ZodObject<{
        selector: z.ZodDiscriminatedUnion<"type", [z.ZodObject<{
            type: z.ZodLiteral<"uuids">;
            values: z.ZodArray<z.ZodString, "many">;
        }, "strip", z.ZodTypeAny, {
            values: string[];
            type: "uuids";
        }, {
            values: string[];
            type: "uuids";
        }>, z.ZodObject<{
            type: z.ZodLiteral<"remarkRegex">;
            pattern: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            type: "remarkRegex";
            pattern: string;
        }, {
            type: "remarkRegex";
            pattern: string;
        }>, z.ZodObject<{
            type: z.ZodLiteral<"tagRegex">;
            pattern: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            type: "tagRegex";
            pattern: string;
        }, {
            type: "tagRegex";
            pattern: string;
        }>, z.ZodObject<{
            type: z.ZodLiteral<"sameTagAsRecipient">;
        }, "strip", z.ZodTypeAny, {
            type: "sameTagAsRecipient";
        }, {
            type: "sameTagAsRecipient";
        }>]>;
        selectFrom: z.ZodOptional<z.ZodEnum<["ALL", "HIDDEN", "NOT_HIDDEN"]>>;
        tagPrefix: z.ZodOptional<z.ZodString>;
        useHostRemarkAsTag: z.ZodOptional<z.ZodBoolean>;
        useHostTagAsTag: z.ZodOptional<z.ZodBoolean>;
    }, "strip", z.ZodTypeAny, {
        selector: {
            values: string[];
            type: "uuids";
        } | {
            type: "remarkRegex";
            pattern: string;
        } | {
            type: "tagRegex";
            pattern: string;
        } | {
            type: "sameTagAsRecipient";
        };
        selectFrom?: "ALL" | "HIDDEN" | "NOT_HIDDEN" | undefined;
        tagPrefix?: string | undefined;
        useHostRemarkAsTag?: boolean | undefined;
        useHostTagAsTag?: boolean | undefined;
    }, {
        selector: {
            values: string[];
            type: "uuids";
        } | {
            type: "remarkRegex";
            pattern: string;
        } | {
            type: "tagRegex";
            pattern: string;
        } | {
            type: "sameTagAsRecipient";
        };
        selectFrom?: "ALL" | "HIDDEN" | "NOT_HIDDEN" | undefined;
        tagPrefix?: string | undefined;
        useHostRemarkAsTag?: boolean | undefined;
        useHostTagAsTag?: boolean | undefined;
    }>, "many">>;
    addVirtualHostAsOutbound: z.ZodOptional<z.ZodBoolean>;
}, "strip", z.ZodTypeAny, {
    injectHosts?: {
        selector: {
            values: string[];
            type: "uuids";
        } | {
            type: "remarkRegex";
            pattern: string;
        } | {
            type: "tagRegex";
            pattern: string;
        } | {
            type: "sameTagAsRecipient";
        };
        selectFrom?: "ALL" | "HIDDEN" | "NOT_HIDDEN" | undefined;
        tagPrefix?: string | undefined;
        useHostRemarkAsTag?: boolean | undefined;
        useHostTagAsTag?: boolean | undefined;
    }[] | undefined;
    addVirtualHostAsOutbound?: boolean | undefined;
}, {
    injectHosts?: {
        selector: {
            values: string[];
            type: "uuids";
        } | {
            type: "remarkRegex";
            pattern: string;
        } | {
            type: "tagRegex";
            pattern: string;
        } | {
            type: "sameTagAsRecipient";
        };
        selectFrom?: "ALL" | "HIDDEN" | "NOT_HIDDEN" | undefined;
        tagPrefix?: string | undefined;
        useHostRemarkAsTag?: boolean | undefined;
        useHostTagAsTag?: boolean | undefined;
    }[] | undefined;
    addVirtualHostAsOutbound?: boolean | undefined;
}>;
export type TCerberusInjector = z.infer<typeof CerberusInjectorSchema>;
export type TCerberusInjectorSelector = z.infer<typeof HostSelectorSchema>;
export type TCerberusInjectorSelectFrom = z.infer<typeof InjectHostsEntrySchema>['selectFrom'];
export {};
//# sourceMappingURL=cerberus-injector.schema.d.ts.map