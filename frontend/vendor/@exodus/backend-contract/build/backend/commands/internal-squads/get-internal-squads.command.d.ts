import { z } from 'zod';
export declare namespace GetInternalSquadsCommand {
    const url: "/api/internal-squads/";
    const TSQ_url: "/api/internal-squads/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            total: z.ZodNumber;
            internalSquads: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                viewPosition: z.ZodInt;
                name: z.ZodString;
                info: z.ZodObject<{
                    membersCount: z.ZodNumber;
                    inboundsCount: z.ZodNumber;
                }, z.core.$strip>;
                inbounds: z.ZodArray<z.ZodObject<{
                    uuid: z.ZodUUID;
                    profileUuid: z.ZodUUID;
                    tag: z.ZodString;
                    type: z.ZodString;
                    network: z.ZodNullable<z.ZodString>;
                    security: z.ZodNullable<z.ZodString>;
                    port: z.ZodNullable<z.ZodNumber>;
                    rawInbound: z.ZodNullable<z.ZodUnknown>;
                }, z.core.$strip>>;
                createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
                updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-internal-squads.command.d.ts.map