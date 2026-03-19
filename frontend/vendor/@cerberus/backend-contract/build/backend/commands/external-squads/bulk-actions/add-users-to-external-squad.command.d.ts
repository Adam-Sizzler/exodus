import { z } from 'zod';
export declare namespace AddUsersToExternalSquadCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../../constants").EndpointDetails;
    const RequestSchema: z.ZodObject<{
        uuid: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        uuid: string;
    }, {
        uuid: string;
    }>;
    type Request = z.infer<typeof RequestSchema>;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            eventSent: z.ZodBoolean;
        }, "strip", z.ZodTypeAny, {
            eventSent: boolean;
        }, {
            eventSent: boolean;
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            eventSent: boolean;
        };
    }, {
        response: {
            eventSent: boolean;
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=add-users-to-external-squad.command.d.ts.map