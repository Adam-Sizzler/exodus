import z from 'zod';
export declare const USER_USAGE_STREAM_MESSAGE_VERSION = "1";
export declare const UserUsageStreamRecordSchema: z.ZodObject<{
    userId: z.ZodString;
    totalBytes: z.ZodString;
}, z.core.$strip>;
export declare const ExodusUserUsageStreamMessageSchema: z.ZodObject<{
    v: z.ZodLiteral<"1">;
    nodeId: z.ZodString;
    ts: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    records: z.ZodPipe<z.ZodPipe<z.ZodString, z.ZodTransform<{
        userId: string;
        totalBytes: string;
    }[], string>>, z.ZodArray<z.ZodObject<{
        userId: z.ZodString;
        totalBytes: z.ZodString;
    }, z.core.$strip>>>;
}, z.core.$strip>;
export type TUserUsageStreamRecord = z.infer<typeof UserUsageStreamRecordSchema>;
export type TExodusUserUsageStreamMessage = z.infer<typeof ExodusUserUsageStreamMessageSchema>;
export declare const SUBSCRIPTION_REQUEST_STREAM_MESSAGE_VERSION = "1";
export declare const ExodusSubscriptionRequestStreamMessageSchema: z.ZodObject<{
    v: z.ZodLiteral<"1">;
    userId: z.ZodString;
    requestAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    requestIp: z.ZodOptional<z.ZodString>;
    userAgent: z.ZodOptional<z.ZodString>;
}, z.core.$strip>;
export type TExodusSubscriptionRequestStreamMessage = z.infer<typeof ExodusSubscriptionRequestStreamMessageSchema>;
export declare const NODE_CONNECTIONS_STREAM_MESSAGE_VERSION = "1";
export declare const NodeConnectionUserSchema: z.ZodObject<{
    userId: z.ZodString;
    ips: z.ZodArray<z.ZodObject<{
        ip: z.ZodString;
        lastSeen: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    }, z.core.$strip>>;
}, z.core.$strip>;
export declare const ExodusNodeConnectionsStreamMessageSchema: z.ZodObject<{
    v: z.ZodLiteral<"1">;
    nodeId: z.ZodString;
    ts: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
    users: z.ZodPipe<z.ZodPipe<z.ZodString, z.ZodTransform<any, string>>, z.ZodArray<z.ZodObject<{
        userId: z.ZodString;
        ips: z.ZodArray<z.ZodObject<{
            ip: z.ZodString;
            lastSeen: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        }, z.core.$strip>>;
    }, z.core.$strip>>>;
}, z.core.$strip>;
export type TNodeConnectionUser = z.infer<typeof NodeConnectionUserSchema>;
export type TExodusNodeConnectionsStreamMessage = z.infer<typeof ExodusNodeConnectionsStreamMessageSchema>;
//# sourceMappingURL=export-stream.schema.d.ts.map