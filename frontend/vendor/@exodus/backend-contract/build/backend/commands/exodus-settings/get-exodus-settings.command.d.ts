import { z } from 'zod';
import { ExodusSettingsSchema } from '../../models';
export declare namespace GetExodusSettingsCommand {
    const url: "/api/exodus-settings";
    const TSQ_url: "/api/exodus-settings";
    const endpointDetails: import("../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: typeof ExodusSettingsSchema;
    }, "strip", z.ZodTypeAny, {
        response: z.infer<typeof ExodusSettingsSchema>;
    }, {
        response: z.infer<typeof ExodusSettingsSchema>;
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
