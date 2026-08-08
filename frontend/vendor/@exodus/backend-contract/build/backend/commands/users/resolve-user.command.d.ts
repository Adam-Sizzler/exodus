import { z } from 'zod';
export declare namespace ResolveUserCommand {
    const url: "/api/users/resolve";
    const TSQ_url: "/api/users/resolve";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestBodySchema: z.ZodObject<{
        id: z.ZodOptional<z.ZodNumber>;
        shortUuid: z.ZodOptional<z.ZodString>;
        username: z.ZodOptional<z.ZodString>;
    }, z.core.$strip>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            id: z.ZodNumber;
            username: z.ZodString;
            shortUuid: z.ZodString;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type RequestBody = z.infer<typeof RequestBodySchema>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=resolve-user.command.d.ts.map