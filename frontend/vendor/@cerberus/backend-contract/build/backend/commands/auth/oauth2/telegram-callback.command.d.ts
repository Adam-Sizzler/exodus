import { z } from 'zod';
export declare namespace TelegramCallbackCommand {
    const url: "/api/auth/oauth2/tg/callback";
    const TSQ_url: "/api/auth/oauth2/tg/callback";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        id: z.ZodNumber;
        first_name: z.ZodString;
        last_name: z.ZodOptional<z.ZodString>;
        username: z.ZodOptional<z.ZodString>;
        photo_url: z.ZodOptional<z.ZodString>;
        auth_date: z.ZodNumber;
        hash: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        id: number;
        first_name: string;
        auth_date: number;
        hash: string;
        username?: string | undefined;
        last_name?: string | undefined;
        photo_url?: string | undefined;
    }, {
        id: number;
        first_name: string;
        auth_date: number;
        hash: string;
        username?: string | undefined;
        last_name?: string | undefined;
        photo_url?: string | undefined;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            accessToken: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            accessToken: string;
        }, {
            accessToken: string;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            accessToken: string;
        };
    }, {
        response: {
            accessToken: string;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=telegram-callback.command.d.ts.map