import { z } from 'zod';
export declare namespace GetUserAccessibleNodesCommand {
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
            userUuid: z.ZodString;
            activeNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                nodeName: z.ZodString;
                countryCode: z.ZodString;
                configProfileUuid: z.ZodString;
                configProfileName: z.ZodString;
                activeSquads: z.ZodArray<z.ZodObject<{
                    squadName: z.ZodString;
                    activeInbounds: z.ZodArray<z.ZodString, "many">;
                }, "strip", z.ZodTypeAny, {
                    activeInbounds: string[];
                    squadName: string;
                }, {
                    activeInbounds: string[];
                    squadName: string;
                }>, "many">;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeSquads: {
                    activeInbounds: string[];
                    squadName: string;
                }[];
                configProfileName: string;
            }, {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeSquads: {
                    activeInbounds: string[];
                    squadName: string;
                }[];
                configProfileName: string;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            userUuid: string;
            activeNodes: {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeSquads: {
                    activeInbounds: string[];
                    squadName: string;
                }[];
                configProfileName: string;
            }[];
        }, {
            userUuid: string;
            activeNodes: {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeSquads: {
                    activeInbounds: string[];
                    squadName: string;
                }[];
                configProfileName: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            userUuid: string;
            activeNodes: {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeSquads: {
                    activeInbounds: string[];
                    squadName: string;
                }[];
                configProfileName: string;
            }[];
        };
    }, {
        response: {
            userUuid: string;
            activeNodes: {
                uuid: string;
                nodeName: string;
                countryCode: string;
                configProfileUuid: string;
                activeSquads: {
                    activeInbounds: string[];
                    squadName: string;
                }[];
                configProfileName: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-user-accessible-nodes.command.d.ts.map