import { z } from 'zod';
export declare namespace UpdateNodeCommand {
    const url: "/api/nodes/";
    const TSQ_url: "/api/nodes/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodOptional<z.ZodString>;
        address: z.ZodOptional<z.ZodString>;
        port: z.ZodOptional<z.ZodNumber>;
        proxyUrl: z.ZodOptional<z.ZodNullable<z.ZodString>>;
        isTrafficTrackingActive: z.ZodOptional<z.ZodBoolean>;
        trafficLimitBytes: z.ZodOptional<z.ZodNumber>;
        notifyPercent: z.ZodOptional<z.ZodNumber>;
        trafficResetDay: z.ZodOptional<z.ZodNumber>;
        countryCode: z.ZodOptional<z.ZodString>;
        consumptionMultiplier: z.ZodOptional<z.ZodPipe<z.ZodNumber, z.ZodTransform<number, number>>>;
        nodeConsumptionMultiplier: z.ZodOptional<z.ZodPipe<z.ZodNumber, z.ZodTransform<number, number>>>;
        configProfile: z.ZodOptional<z.ZodObject<{
            activeConfigProfileUuid: z.ZodUUID;
            activeInbounds: z.ZodArray<z.ZodUUID>;
        }, z.core.$strip>>;
        providerUuid: z.ZodOptional<z.ZodNullable<z.ZodUUID>>;
        tags: z.ZodOptional<z.ZodArray<z.ZodString>>;
        activePluginUuid: z.ZodOptional<z.ZodNullable<z.ZodUUID>>;
        note: z.ZodOptional<z.ZodNullable<z.ZodString>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            id: z.ZodNumber;
            name: z.ZodString;
            address: z.ZodString;
            port: z.ZodNullable<z.ZodInt>;
            proxyUrl: z.ZodNullable<z.ZodString>;
            isConnected: z.ZodBoolean;
            isDisabled: z.ZodBoolean;
            isConnecting: z.ZodBoolean;
            lastStatusChange: z.ZodNullable<z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>>;
            lastStatusMessage: z.ZodNullable<z.ZodString>;
            isTrafficTrackingActive: z.ZodBoolean;
            trafficResetDay: z.ZodNullable<z.ZodInt>;
            trafficLimitBytes: z.ZodNullable<z.ZodNumber>;
            trafficUsedBytes: z.ZodNullable<z.ZodNumber>;
            notifyPercent: z.ZodNullable<z.ZodInt>;
            viewPosition: z.ZodInt;
            countryCode: z.ZodString;
            consumptionMultiplier: z.ZodNumber;
            nodeConsumptionMultiplier: z.ZodNumber;
            tags: z.ZodArray<z.ZodString>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            configProfile: z.ZodObject<{
                activeConfigProfileUuid: z.ZodNullable<z.ZodUUID>;
                activeInbounds: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodUUID;
                    profileUuid: z.ZodUUID;
                    tag: z.ZodString;
                    type: z.ZodString;
                    network: z.ZodNullable<z.ZodString>;
                    security: z.ZodNullable<z.ZodString>;
                    port: z.ZodNullable<z.ZodNumber>;
                    rawInbound: z.ZodNullable<z.ZodUnknown>;
                }, z.core.$strip>>;
            }, z.core.$strip>;
            providerUuid: z.ZodNullable<z.ZodUUID>;
            provider: z.ZodNullable<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
                faviconLink: z.ZodNullable<z.ZodString>;
                loginUrl: z.ZodNullable<z.ZodString>;
                createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            }, z.core.$strip>>;
            activePluginUuid: z.ZodNullable<z.ZodUUID>;
            system: z.ZodNullable<z.ZodObject<{
                info: z.ZodObject<{
                    arch: z.ZodString;
                    cpus: z.ZodInt;
                    cpuModel: z.ZodString;
                    memoryTotal: z.ZodNumber;
                    hostname: z.ZodString;
                    platform: z.ZodString;
                    release: z.ZodString;
                    type: z.ZodString;
                    version: z.ZodString;
                    networkInterfaces: z.ZodArray<z.ZodString>;
                }, z.core.$strip>;
                stats: z.ZodObject<{
                    memoryFree: z.ZodNumber;
                    memoryUsed: z.ZodNumber;
                    uptime: z.ZodNumber;
                    loadAvg: z.ZodArray<z.ZodNumber>;
                    interface: z.ZodNullable<z.ZodObject<{
                        interface: z.ZodString;
                        rxBytesPerSec: z.ZodNumber;
                        txBytesPerSec: z.ZodNumber;
                        rxTotal: z.ZodNumber;
                        txTotal: z.ZodNumber;
                    }, z.core.$strip>>;
                }, z.core.$strip>;
            }, z.core.$strip>>;
            versions: z.ZodNullable<z.ZodObject<{
                singbox: z.ZodString;
                node: z.ZodString;
            }, z.core.$strip>>;
            singboxUptime: z.ZodNumber;
            usersOnline: z.ZodNumber;
            note: z.ZodNullable<z.ZodString>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Request = RequestBody;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update.command.d.ts.map