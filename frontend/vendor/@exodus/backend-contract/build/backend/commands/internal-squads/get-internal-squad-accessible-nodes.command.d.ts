import { z } from 'zod';
export declare namespace GetInternalSquadAccessibleNodesCommand {
    const url: (uuid: string) => string;
    const TSQ_url: string;
    const endpointDetails: import("../../constants").EndpointDetails;
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
            squadUuid: z.ZodString;
            accessibleNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                nodeName: z.ZodString;
                countryCode: z.ZodString;
                configProfileUuid: z.ZodString;
                configProfileName: z.ZodString;
                activeInbounds: z.ZodArray<z.ZodString, "many">;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeInbounds: string[];
                configProfileName: string;
            }, {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeInbounds: string[];
                configProfileName: string;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            squadUuid: string;
            accessibleNodes: {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeInbounds: string[];
                configProfileName: string;
            }[];
        }, {
            squadUuid: string;
            accessibleNodes: {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeInbounds: string[];
                configProfileName: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            squadUuid: string;
            accessibleNodes: {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeInbounds: string[];
                configProfileName: string;
            }[];
        };
    }, {
        response: {
            squadUuid: string;
            accessibleNodes: {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeInbounds: string[];
                configProfileName: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-internal-squad-accessible-nodes.command.d.ts.map