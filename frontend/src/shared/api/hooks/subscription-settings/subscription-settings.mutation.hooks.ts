import { notifications } from '@mantine/notifications'
import { UpdateSubscriptionSettingsCommand } from '@exodus/backend-contract'

import { createMutationHook } from '../../tsq-helpers'

export const useUpdateSubscriptionSettings = createMutationHook({
    endpoint: UpdateSubscriptionSettingsCommand.TSQ_url,
    bodySchema: UpdateSubscriptionSettingsCommand.RequestBodySchema,
    responseSchema: UpdateSubscriptionSettingsCommand.ResponseSchema,
    requestMethod: UpdateSubscriptionSettingsCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Subscription settings updated successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: 'Update Subscription Settings',
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})
