import { notifications } from '@mantine/notifications'
import {
    CreateApiTokenCommand,
    DeleteApiTokenCommand,
    GetOttCommand
} from '@exodus/backend-contract'

import { createMutationHook } from '@shared/api/tsq-helpers/create-mutation-hook'

export const useCreateApiToken = createMutationHook({
    endpoint: CreateApiTokenCommand.TSQ_url,
    bodySchema: CreateApiTokenCommand.RequestBodySchema,
    responseSchema: CreateApiTokenCommand.ResponseSchema,
    requestMethod: CreateApiTokenCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Api token created successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Create Api Token`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useDeleteApiToken = createMutationHook({
    endpoint: DeleteApiTokenCommand.TSQ_url,
    routeParamsSchema: DeleteApiTokenCommand.RequestParamSchema,
    requestMethod: DeleteApiTokenCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Api token deleted successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Delete Api Token`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useIssueOtt = createMutationHook({
    endpoint: GetOttCommand.TSQ_url,
    responseSchema: GetOttCommand.ResponseSchema,
    requestMethod: GetOttCommand.endpointDetails.REQUEST_METHOD
})
