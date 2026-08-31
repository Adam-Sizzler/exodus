import { z } from 'zod';
export declare namespace GetInternalSquadCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestParamSchema: z.ZodObject<{
        uuid: z.ZodUUID;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            viewPosition: z.ZodInt;
            name: z.ZodString;
            tags: z.ZodDefault<z.ZodArray<z.ZodString>>;
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
    type RequestParam = z.infer<typeof RequestParamSchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-internal-squad.command.d.ts.map