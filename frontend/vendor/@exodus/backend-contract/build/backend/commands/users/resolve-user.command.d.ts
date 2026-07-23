import { z } from 'zod';
export declare namespace ResolveUserCommand {
    const url: "/api/users/resolve";
    const TSQ_url: "/api/users/resolve";
    const endpointDetails: import("../../constants").EndpointDetails;
    const RequestSchema: z.ZodEffects<z.ZodObject<{
        uuid: z.ZodOptional<z.ZodString>;
        id: z.ZodOptional<z.ZodNumber>;
        shortUuid: z.ZodOptional<z.ZodString>;
        username: z.ZodOptional<z.ZodString>;
    }, "strip", z.ZodTypeAny, {
        uuid?: string | undefined;
        username?: string | undefined;
        id?: number | undefined;
        shortUuid?: string | undefined;
    }, {
        uuid?: string | undefined;
        username?: string | undefined;
        id?: number | undefined;
        shortUuid?: string | undefined;
    }>, {
        uuid?: string | undefined;
        username?: string | undefined;
        id?: number | undefined;
        shortUuid?: string | undefined;
    }, {
        uuid?: string | undefined;
        username?: string | undefined;
        id?: number | undefined;
        shortUuid?: string | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            uuid: z.ZodString;
            username: z.ZodString;
            id: z.ZodNumber;
            shortUuid: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            uuid: string;
            username: string;
            id: number;
            shortUuid: string;
        }, {
            uuid: string;
            username: string;
            id: number;
            shortUuid: string;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            uuid: string;
            username: string;
            id: number;
            shortUuid: string;
        };
    }, {
        response: {
            uuid: string;
            username: string;
            id: number;
            shortUuid: string;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=resolve-user.command.d.ts.map