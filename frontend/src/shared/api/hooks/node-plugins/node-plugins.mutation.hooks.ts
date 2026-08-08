import { notifications } from '@mantine/notifications'
import {
    CloneNodePluginCommand,
    CreateNodePluginCommand,
    DeleteNodePluginCommand,
    PluginExecutorCommand,
    ReorderNodePluginCommand,
    UpdateNodePluginCommand
} from '@exodus/backend-contract'

import { createMutationHook } from '../../tsq-helpers'

export const useUpdateNodePlugin = createMutationHook({
    endpoint: UpdateNodePluginCommand.TSQ_url,
    bodySchema: UpdateNodePluginCommand.RequestBodySchema,
    responseSchema: UpdateNodePluginCommand.ResponseSchema,
    requestMethod: UpdateNodePluginCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node plugin updated successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Update Node Plugin`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useCreateNodePlugin = createMutationHook({
    endpoint: CreateNodePluginCommand.TSQ_url,
    bodySchema: CreateNodePluginCommand.RequestBodySchema,
    responseSchema: CreateNodePluginCommand.ResponseSchema,
    requestMethod: CreateNodePluginCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node plugin created successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Create Node Plugin`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useDeleteNodePlugin = createMutationHook({
    endpoint: DeleteNodePluginCommand.TSQ_url,
    routeParamsSchema: DeleteNodePluginCommand.RequestParamSchema,
    requestMethod: DeleteNodePluginCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node plugin deleted successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Delete Node Plugin`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useReorderNodePlugins = createMutationHook({
    endpoint: ReorderNodePluginCommand.TSQ_url,
    bodySchema: ReorderNodePluginCommand.RequestBodySchema,
    responseSchema: ReorderNodePluginCommand.ResponseSchema,
    requestMethod: ReorderNodePluginCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onError: (error) => {
            notifications.show({
                title: `Reorder Node Plugins`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useCloneNodePlugin = createMutationHook({
    endpoint: CloneNodePluginCommand.TSQ_url,
    bodySchema: CloneNodePluginCommand.RequestBodySchema,
    responseSchema: CloneNodePluginCommand.ResponseSchema,
    requestMethod: CloneNodePluginCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node plugin cloned successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Clone Node Plugin`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useNodePluginExecutor = createMutationHook({
    endpoint: PluginExecutorCommand.TSQ_url,
    bodySchema: PluginExecutorCommand.RequestBodySchema,
    requestMethod: PluginExecutorCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Request sent',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Node Plugin Executor`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})
