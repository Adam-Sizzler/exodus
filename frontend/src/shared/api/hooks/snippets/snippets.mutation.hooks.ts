import { notifications } from '@mantine/notifications'
import {
    CreateSnippetCommand,
    DeleteSnippetCommand,
    UpdateSnippetCommand
} from '@exodus/backend-contract'

import { createMutationHook } from '../../tsq-helpers'

export const useUpdateSnippet = createMutationHook({
    endpoint: UpdateSnippetCommand.TSQ_url,
    bodySchema: UpdateSnippetCommand.RequestBodySchema,
    responseSchema: UpdateSnippetCommand.ResponseSchema,
    requestMethod: UpdateSnippetCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Snippet updated successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Update Snippet`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useDeleteSnippet = createMutationHook({
    endpoint: DeleteSnippetCommand.TSQ_url,
    bodySchema: DeleteSnippetCommand.RequestBodySchema,
    requestMethod: DeleteSnippetCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Snippet deleted successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Delete Snippet`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useCreateSnippet = createMutationHook({
    endpoint: CreateSnippetCommand.TSQ_url,
    responseSchema: CreateSnippetCommand.ResponseSchema,
    bodySchema: CreateSnippetCommand.RequestBodySchema,
    requestMethod: CreateSnippetCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Snippet created successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Create Snippet`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})
