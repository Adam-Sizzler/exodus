import {
    BulkNodesActionsCommand,
    BulkNodesProfileModificationCommand,
    CreateNodeCommand,
    DeleteNodeCommand,
    DisableNodeCommand,
    EnableNodeCommand,
    ReorderNodeCommand,
    ResetNodeTrafficCommand,
    RestartAllNodesCommand,
    RestartNodeCommand,
    UpdateNodeCommand
} from '@exodus/backend-contract'
import { notifications } from '@mantine/notifications'
import { z } from 'zod'

import { createMutationHook } from '../../tsq-helpers'

const nodeConnectionSchema = {
    apiSchema: z.enum(['mtls', 'tls']).default('mtls'),
    apiPath: z.string().trim().min(1).default('/'),
    grpcAuthToken: z.string().trim().length(64).optional(),
    activePluginUuid: z.string().uuid().nullable().optional()
}

export const createNodeRequestSchema = CreateNodeCommand.RequestSchema.extend(nodeConnectionSchema)
export type CreateNodeRequest = z.infer<typeof createNodeRequestSchema>

export const updateNodeRequestSchema = UpdateNodeCommand.RequestSchema.extend({
    apiSchema: nodeConnectionSchema.apiSchema.optional(),
    apiPath: nodeConnectionSchema.apiPath.optional(),
    grpcAuthToken: nodeConnectionSchema.grpcAuthToken,
    activePluginUuid: nodeConnectionSchema.activePluginUuid
})
export const updateNodeFormSchema = updateNodeRequestSchema.omit({ uuid: true })
export type UpdateNodeRequest = z.infer<typeof updateNodeRequestSchema>

const pluginConfigSchema = z
    .record(z.unknown())
    .refine((value) => !('torrentBlocker' in value), {
        message: 'torrentBlocker plugin is not supported'
    })
    .refine((value) => !('connectionDrop' in value), {
        message: 'connectionDrop plugin is not supported'
    })

const nodePluginResponseSchema = z.object({
    response: z.object({
        uuid: z.string(),
        name: z.string(),
        pluginConfig: pluginConfigSchema,
        viewPosition: z.number().optional().default(0)
    })
})

const nodePluginsListResponseSchema = z.object({
    response: z.object({
        nodePlugins: z.array(nodePluginResponseSchema.shape.response),
        total: z.number().optional().default(0)
    })
})

const createNodePluginRequestSchema = z.object({
    name: z.string().trim().min(1),
    pluginConfig: pluginConfigSchema.optional()
})

const updateNodePluginRequestSchema = z.object({
    name: z.string().trim().min(1).optional(),
    pluginConfig: pluginConfigSchema.optional(),
    viewPosition: z.number().optional()
})

const deleteNodePluginResponseSchema = z.object({
    response: z.object({
        isDeleted: z.boolean()
    })
})

const reorderNodePluginsRequestSchema = z.object({
    items: z.array(
        z.object({
            uuid: z.string().uuid(),
            viewPosition: z.number()
        })
    )
})

const cloneNodePluginRequestSchema = z.object({
    cloneFromUuid: z.string().uuid(),
    name: z.string().trim().min(1).optional()
})

const nodePluginExecutorRequestSchema = z.object({
    command: z.discriminatedUnion('command', [
        z.object({
            command: z.literal('blockIps'),
            ips: z.array(
                z.object({
                    ip: z.string().min(1),
                    timeout: z.number().int().min(0)
                })
            )
        }),
        z.object({
            command: z.literal('unblockIps'),
            ips: z.array(z.string().min(1))
        }),
        z.object({
            command: z.literal('recreateTables')
        })
    ]),
    targetNodes: z.object({
        target: z.literal('specificNodes'),
        nodeUuids: z.array(z.string().uuid()).min(1)
    })
})

const nodePluginExecutorResponseSchema = z.object({
    response: z.object({
        eventSent: z.boolean()
    })
})

export const useCreateNode = createMutationHook({
    endpoint: CreateNodeCommand.TSQ_url,
    bodySchema: createNodeRequestSchema,
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

export const useCreateNodePlugin = createMutationHook({
    endpoint: '/api/node-plugins',
    bodySchema: createNodePluginRequestSchema,
    responseSchema: nodePluginResponseSchema,
    requestMethod: 'post',
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node plugin created successfully',
                color: 'teal'
            })
        }
    }
})

export const useUpdateNodePlugin = createMutationHook({
    endpoint: '/api/node-plugins/:uuid',
    bodySchema: updateNodePluginRequestSchema,
    responseSchema: nodePluginResponseSchema,
    routeParamsSchema: z.object({ uuid: z.string().uuid() }),
    requestMethod: 'patch',
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node plugin updated successfully',
                color: 'teal'
            })
        }
    }
})

export const useDeleteNodePlugin = createMutationHook({
    endpoint: '/api/node-plugins/:uuid',
    responseSchema: deleteNodePluginResponseSchema,
    routeParamsSchema: z.object({ uuid: z.string().uuid() }),
    requestMethod: 'delete',
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node plugin deleted successfully',
                color: 'teal'
            })
        }
    }
})

export const useReorderNodePlugins = createMutationHook({
    endpoint: '/api/node-plugins/actions/reorder',
    bodySchema: reorderNodePluginsRequestSchema,
    responseSchema: nodePluginsListResponseSchema,
    requestMethod: 'post'
})

export const useCloneNodePlugin = createMutationHook({
    endpoint: '/api/node-plugins/actions/clone',
    bodySchema: cloneNodePluginRequestSchema,
    responseSchema: nodePluginResponseSchema,
    requestMethod: 'post',
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Node plugin cloned successfully',
                color: 'teal'
            })
        }
    }
})

export const useNodePluginExecutor = createMutationHook({
    endpoint: '/api/node-plugins/executor',
    bodySchema: nodePluginExecutorRequestSchema,
    responseSchema: nodePluginExecutorResponseSchema,
    requestMethod: 'post',
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

export const useUpdateNode = createMutationHook({
    endpoint: UpdateNodeCommand.TSQ_url,
    bodySchema: updateNodeRequestSchema,
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
    responseSchema: DeleteNodeCommand.ResponseSchema,
    routeParamsSchema: DeleteNodeCommand.RequestSchema,
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
    routeParamsSchema: EnableNodeCommand.RequestSchema,
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
    routeParamsSchema: DisableNodeCommand.RequestSchema,
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
    responseSchema: RestartAllNodesCommand.ResponseSchema,
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
    endpoint: ReorderNodeCommand.TSQ_url,
    bodySchema: ReorderNodeCommand.RequestSchema,
    responseSchema: ReorderNodeCommand.ResponseSchema,
    requestMethod: ReorderNodeCommand.endpointDetails.REQUEST_METHOD,
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
    responseSchema: RestartNodeCommand.ResponseSchema,
    routeParamsSchema: RestartNodeCommand.RequestSchema,
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
    responseSchema: ResetNodeTrafficCommand.ResponseSchema,
    routeParamsSchema: ResetNodeTrafficCommand.RequestSchema,
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
    responseSchema: BulkNodesProfileModificationCommand.ResponseSchema,
    bodySchema: BulkNodesProfileModificationCommand.RequestSchema,
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
    responseSchema: BulkNodesActionsCommand.ResponseSchema,
    bodySchema: BulkNodesActionsCommand.RequestSchema,
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
