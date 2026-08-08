import { z } from 'zod';
export declare namespace UpdateConfigProfileCommand {
    const url: "/api/config-profiles/";
    const TSQ_url: "/api/config-profiles/";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        uuid: z.ZodUUID;
        name: z.ZodOptional<z.ZodString>;
        config: z.ZodOptional<z.ZodObject<{}, z.core.$loose>>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodUUID;
            viewPosition: z.ZodInt;
            name: z.ZodString;
            config: z.ZodUnknown;
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
            nodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                name: z.ZodString;
                countryCode: z.ZodString;
            }, z.core.$strip>>;
            createdAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
            updatedAt: z.ZodPipe<z.ZodISODateTime, z.ZodTransform<Date, string>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=update-config-profile.command.d.ts.map