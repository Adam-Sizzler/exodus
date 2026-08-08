import { z } from 'zod';
export declare namespace UpdateInternalSquadCommand {
    const url: "/api/internal-squads/";
    const TSQ_url: "/api/internal-squads/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodOptional<z.ZodString>;
        inbounds: z.ZodOptional<z.ZodArray<z.ZodUUID>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
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
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-internal-squad.command.d.ts.map