import { UpdateCerberusSettingsCommand } from '@cerberus/backend-contract'
import { notifications } from '@mantine/notifications'

import { createMutationHook } from '../../tsq-helpers'

export const useUpdateCerberusSettings = createMutationHook({
    endpoint: UpdateCerberusSettingsCommand.TSQ_url,
    bodySchema: UpdateCerberusSettingsCommand.RequestSchema,
    responseSchema: UpdateCerberusSettingsCommand.ResponseSchema,
    requestMethod: UpdateCerberusSettingsCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Cerberus settings updated successfully',
                color: 'teal'
            })
        }
    }
})
