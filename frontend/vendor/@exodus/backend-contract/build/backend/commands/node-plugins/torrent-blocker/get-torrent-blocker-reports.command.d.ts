import { z } from 'zod';
export declare namespace GetTorrentBlockerReportsCommand {
    const url: "/api/node-plugins/torrent-blocker";
    const TSQ_url: "/api/node-plugins/torrent-blocker";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestQuerySchema: z.ZodObject<{
        start: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
        size: z.ZodDefault<z.ZodCoercedNumber<unknown>>;
        filters: z.ZodOptional<z.ZodPreprocess<z.ZodArray<z.ZodObject<{
            id: z.ZodString;
            value: z.ZodUnknown;
        }, z.core.$strip>>>>;
        filterModes: z.ZodOptional<z.ZodPreprocess<z.ZodRecord<z.ZodString, z.ZodString>>>;
        globalFilterMode: z.ZodOptional<z.ZodString>;
        sorting: z.ZodOptional<z.ZodPreprocess<z.ZodArray<z.ZodObject<{
            id: z.ZodString;
            desc: z.ZodBoolean;
        }, z.core.$strip>>>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            records: z.ZodArray<z.ZodObject<{
                id: z.ZodNumber;
                userId: z.ZodNumber;
                nodeId: z.ZodNumber;
                user: z.ZodObject<{
                    username: z.ZodString;
                }, z.core.$strip>;
                node: z.ZodObject<{
                    uuid: z.ZodUUID;
                    name: z.ZodString;
                    countryCode: z.ZodString;
                }, z.core.$strip>;
                report: z.ZodObject<{
                    actionReport: z.ZodObject<{
                        blocked: z.ZodBoolean;
                        ip: z.ZodString;
                        blockDuration: z.ZodNumber;
                        willUnblockAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                        userId: z.ZodString;
                        processedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                    }, z.core.$strip>;
                    xrayReport: z.ZodObject<{
                        email: z.ZodNullable<z.ZodString>;
                        level: z.ZodNullable<z.ZodNumber>;
                        protocol: z.ZodNullable<z.ZodString>;
                        network: z.ZodString;
                        source: z.ZodNullable<z.ZodString>;
                        destination: z.ZodString;
                        routeTarget: z.ZodNullable<z.ZodString>;
                        originalTarget: z.ZodNullable<z.ZodString>;
                        inboundTag: z.ZodNullable<z.ZodString>;
                        inboundName: z.ZodNullable<z.ZodString>;
                        inboundLocal: z.ZodNullable<z.ZodString>;
                        outboundTag: z.ZodNullable<z.ZodString>;
                        ts: z.ZodNumber;
                    }, z.core.$strip>;
                }, z.core.$strip>;
                createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            }, z.core.$strip>>;
            total: z.ZodNumber;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestQuery = z.infer<typeof RequestQuerySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-torrent-blocker-reports.command.d.ts.map