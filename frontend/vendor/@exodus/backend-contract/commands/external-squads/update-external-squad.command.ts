import { z } from 'zod';

import { EXTERNAL_SQUADS_ROUTES, REST_API } from '../../api';
import { getEndpointDetails, SUBSCRIPTION_TEMPLATE_TYPE } from '../../constants';
import {
    CustomRemarksSchema,
    ExternalSquadHostOverridesSchema,
    ExternalSquadResponseHeadersAddSchema,
    ExternalSquadResponseHeadersRemoveSchema,
    ExternalSquadSchema,
    ExternalSquadSubscriptionSettingsSchema,
    HwidSettingsSchema,
} from '../../models';

export namespace UpdateExternalSquadCommand {
    export const url = REST_API.EXTERNAL_SQUADS.UPDATE;
    export const TSQ_url = url;

    export const endpointDetails = getEndpointDetails(
        EXTERNAL_SQUADS_ROUTES.UPDATE,
        'patch',
        'Update external squad',
        { scope: 'update', kind: 'write' },
    );

    export const RequestBodySchema = z.object({
        uuid: z.uuid().describe('UUID of the external squad'),
        name: z
            .string()
            .min(2)
            .max(30)
            .regex(
                /^[A-Za-z0-9_\s-]+$/,
                'Name can only contain letters, numbers, underscores, dashes and spaces',
            )
            .optional(),
        templates: z
            .array(
                z.object({
                    templateUuid: z.uuid().describe('UUID of the subscription template'),
                    templateType: z
                        .enum(SUBSCRIPTION_TEMPLATE_TYPE)
                        .describe('Type of the subscription template'),
                }),
            )
            .optional(),
        subscriptionSettings: ExternalSquadSubscriptionSettingsSchema.optional(),
        hostOverrides: ExternalSquadHostOverridesSchema.optional(),
        responseHeadersAdd: ExternalSquadResponseHeadersAddSchema.optional(),
        responseHeadersRemove: ExternalSquadResponseHeadersRemoveSchema.optional(),
        hwidSettings: HwidSettingsSchema.nullish(),
        customRemarks: CustomRemarksSchema.nullish(),
        subpageConfigUuid: z.uuid().nullish(),
    });

    export const ResponseSchema = z.object({
        response: ExternalSquadSchema,
    });

    export type RequestBody = z.infer<typeof RequestBodySchema>;
    export type Response = z.infer<typeof ResponseSchema>;
}
