import { UpdateExodusSettingsCommand } from '@exodus/backend-contract'
import { notifications } from '@mantine/notifications'

import { createMutationHook } from '../../tsq-helpers'

export const useUpdateExodusSettings = createMutationHook({
    endpoint: UpdateExodusSettingsCommand.TSQ_url,
    bodySchema: UpdateExodusSettingsCommand.RequestSchema,
    responseSchema: UpdateExodusSettingsCommand.ResponseSchema,
    requestMethod: UpdateExodusSettingsCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Exodus settings updated successfully',
                color: 'teal'
            })
        }
    }
})
