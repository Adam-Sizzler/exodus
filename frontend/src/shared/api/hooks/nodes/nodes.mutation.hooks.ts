import { z } from 'zod'
import { notifications } from '@mantine/notifications'
import {
    BulkNodesActionsCommand,
    BulkNodesProfileModificationCommand,
    BulkNodesUpdateCommand,
    CreateNodeCommand,
    DeleteNodeCommand,
    DisableNodeCommand,
    EnableNodeCommand,
    ReorderNodesCommand,
    ResetNodeTrafficCommand,
    RestartAllNodesCommand,
    RestartNodeCommand,
    UpdateNodeCommand
} from '@exodus/backend-contract'

import { createMutationHook } from '../../tsq-helpers'

export const useCreateNode = createMutationHook({
    endpoint: CreateNodeCommand.TSQ_url,
    bodySchema: CreateNodeCommand.RequestBodySchema,
    responseSchema: CreateNodeCommand.ResponseSchema,
    requestMethod: CreateNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node created successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Create Node`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useUpdateNode = createMutationHook({
    endpoint: UpdateNodeCommand.TSQ_url,
    bodySchema: UpdateNodeCommand.RequestBodySchema.extend({
        apiSchema: z.enum(['mtls', 'tls']).optional(),
        apiPath: z.string().optional(),
        grpcAuthToken: z.string().optional()
    }),
    responseSchema: UpdateNodeCommand.ResponseSchema,
    requestMethod: UpdateNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node updated successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Update Node`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useDeleteNode = createMutationHook({
    endpoint: DeleteNodeCommand.TSQ_url,
    routeParamsSchema: DeleteNodeCommand.RequestParamSchema,
    requestMethod: DeleteNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node deleted successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Delete Node`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useEnableNode = createMutationHook({
    endpoint: EnableNodeCommand.TSQ_url,
    responseSchema: EnableNodeCommand.ResponseSchema,
    routeParamsSchema: EnableNodeCommand.RequestParamSchema,
    requestMethod: EnableNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node enabled successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Enable Node`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useDisableNode = createMutationHook({
    endpoint: DisableNodeCommand.TSQ_url,
    responseSchema: DisableNodeCommand.ResponseSchema,
    routeParamsSchema: DisableNodeCommand.RequestParamSchema,
    requestMethod: DisableNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node disabled successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Disable Node`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useRestartAllNodes = createMutationHook({
    endpoint: RestartAllNodesCommand.TSQ_url,
    bodySchema: RestartAllNodesCommand.RequestBodySchema,
    requestMethod: RestartAllNodesCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Please wait for the nodes to reconnect',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Restart All Nodes`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})
export const useReorderNodes = createMutationHook({
    endpoint: ReorderNodesCommand.TSQ_url,
    bodySchema: ReorderNodesCommand.RequestBodySchema,
    responseSchema: ReorderNodesCommand.ResponseSchema,
    requestMethod: ReorderNodesCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onError: (error) => {
            notifications.show({
                title: `Reorder Nodes`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useRestartNode = createMutationHook({
    endpoint: RestartNodeCommand.TSQ_url,
    routeParamsSchema: RestartNodeCommand.RequestParamSchema,
    bodySchema: RestartNodeCommand.RequestBodySchema,
    requestMethod: RestartNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node restarted successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Restart Node`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useResetNodeTraffic = createMutationHook({
    endpoint: ResetNodeTrafficCommand.TSQ_url,
    routeParamsSchema: ResetNodeTrafficCommand.RequestParamSchema,
    requestMethod: ResetNodeTrafficCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node traffic reset successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Reset Node Traffic`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useBulkNodesProfileModification = createMutationHook({
    endpoint: BulkNodesProfileModificationCommand.TSQ_url,
    bodySchema: BulkNodesProfileModificationCommand.RequestBodySchema,
    requestMethod: BulkNodesProfileModificationCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Task added to queue successfully.',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Bulk Nodes Profile Modification`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useBulkNodesActions = createMutationHook({
    endpoint: BulkNodesActionsCommand.TSQ_url,
    bodySchema: BulkNodesActionsCommand.RequestBodySchema,
    requestMethod: BulkNodesActionsCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Actions added to queue successfully.',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Bulk Nodes Actions`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useBulkNodesUpdate = createMutationHook({
    endpoint: BulkNodesUpdateCommand.TSQ_url,
    bodySchema: BulkNodesUpdateCommand.RequestBodySchema,
    requestMethod: BulkNodesUpdateCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Nodes updated successfully.',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Bulk Nodes Update`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})
