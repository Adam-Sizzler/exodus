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

const SUBSCRIPTION_CONNECTIONS_API = {
    CREATE: '/api/subscription-connections',
    UPDATE: '/api/subscription-connections',
    DELETE: '/api/subscription-connections/:uuid',
    ENABLE: '/api/subscription-connections/:uuid/actions/enable',
    DISABLE: '/api/subscription-connections/:uuid/actions/disable',
    RESTART_ALL: '/api/subscription-connections/actions/restart-all',
    REORDER: '/api/subscription-connections/actions/reorder',
    RESTART: '/api/subscription-connections/:uuid/actions/restart',
    RESET_TRAFFIC: '/api/subscription-connections/:uuid/actions/reset-traffic',
    BULK_PROFILE_MODIFICATION: '/api/subscription-connections/bulk-actions/profile-modification',
    BULK_ACTIONS: '/api/subscription-connections/bulk-actions'
} as const

const TAG_REGEX = /^[A-Z0-9_:]+$/

const isValidPublicDomain = (value: string) => {
    const trimmed = value.trim()
    if (trimmed === '') {
        return true
    }

    const candidate = trimmed.includes('://') ? trimmed : `https://${trimmed}`

    try {
        const parsed = new URL(candidate)
        if (!parsed.hostname) {
            return false
        }
        if (parsed.pathname && parsed.pathname !== '/') {
            return false
        }
        if (parsed.search || parsed.hash) {
            return false
        }
        if (parsed.username || parsed.password) {
            return false
        }
        return true
    } catch {
        return false
    }
}

const publicDomainSchema = z
    .string()
    .trim()
    .max(255)
    .refine(isValidPublicDomain, {
        message: 'Invalid public domain'
    })
    .nullable()
    .optional()

const createSubscriptionConnectionSchemaBase = z.object({
    name: z.string().min(3).max(30),
    address: z.string().min(2),
    publicDomain: publicDomainSchema,
    port: z.number().int().min(1).max(65535).optional(),
    apiSchema: z.enum(['mtls', 'tls']).default('mtls'),
    apiPath: z.string().trim().min(1).default('/'),
    grpcAuthToken: z.string().trim().length(64).optional(),
    subpageConfigUuid: z.string().uuid().nullable().optional(),
    providerUuid: z.string().uuid().nullable().optional(),
    tags: z.array(z.string().regex(TAG_REGEX).max(36)).max(10).optional()
})

export const createSubscriptionConnectionSchema = createSubscriptionConnectionSchemaBase

export type CreateSubscriptionConnectionRequest = z.infer<typeof createSubscriptionConnectionSchema>

const updateSubscriptionConnectionSchemaBase = z.object({
    uuid: z.string().uuid(),
    name: z.string().min(3).max(30).optional(),
    address: z.string().min(2).optional(),
    publicDomain: publicDomainSchema,
    port: z.number().int().min(1).max(65535).optional(),
    apiSchema: z.enum(['mtls', 'tls']).optional(),
    apiPath: z.string().trim().min(1).optional(),
    grpcAuthToken: z.string().trim().length(64).optional(),
    subpageConfigUuid: z.string().uuid().nullable().optional(),
    providerUuid: z.string().uuid().nullable().optional(),
    tags: z.array(z.string().regex(TAG_REGEX).max(36)).max(10).optional()
})

export const updateSubscriptionConnectionSchema = updateSubscriptionConnectionSchemaBase

export const updateSubscriptionConnectionFormSchema = updateSubscriptionConnectionSchemaBase
    .omit({ uuid: true })

export type UpdateSubscriptionConnectionRequest = z.infer<typeof updateSubscriptionConnectionSchema>

export const useCreateSubscriptionConnection = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.CREATE,
    bodySchema: createSubscriptionConnectionSchema,
    responseSchema: CreateNodeCommand.ResponseSchema,
    requestMethod: CreateNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Connection created successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Create Connection`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useUpdateSubscriptionConnection = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.UPDATE,
    bodySchema: updateSubscriptionConnectionSchema,
    responseSchema: UpdateNodeCommand.ResponseSchema,
    requestMethod: UpdateNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Connection updated successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Update Connection`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useDeleteSubscriptionConnection = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.DELETE,
    responseSchema: DeleteNodeCommand.ResponseSchema,
    routeParamsSchema: DeleteNodeCommand.RequestSchema,
    requestMethod: DeleteNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Connection deleted successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Delete Connection`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useEnableSubscriptionConnection = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.ENABLE,
    responseSchema: EnableNodeCommand.ResponseSchema,
    routeParamsSchema: EnableNodeCommand.RequestSchema,
    requestMethod: EnableNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Connection enabled successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Enable Connection`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useDisableSubscriptionConnection = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.DISABLE,
    responseSchema: DisableNodeCommand.ResponseSchema,
    routeParamsSchema: DisableNodeCommand.RequestSchema,
    requestMethod: DisableNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Connection disabled successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Disable Connection`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useRestartAllSubscriptionConnections = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.RESTART_ALL,
    responseSchema: RestartAllNodesCommand.ResponseSchema,
    bodySchema: RestartAllNodesCommand.RequestBodySchema,
    requestMethod: RestartAllNodesCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Please wait for the connections to reconnect',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Restart All Connections`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useReorderSubscriptionConnections = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.REORDER,
    bodySchema: ReorderNodeCommand.RequestSchema,
    responseSchema: ReorderNodeCommand.ResponseSchema,
    requestMethod: ReorderNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onError: (error) => {
            notifications.show({
                title: `Reorder Connections`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useRestartSubscriptionConnection = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.RESTART,
    responseSchema: RestartNodeCommand.ResponseSchema,
    routeParamsSchema: RestartNodeCommand.RequestSchema,
    requestMethod: RestartNodeCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Connection restarted successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Restart Connection`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useResetSubscriptionConnectionTraffic = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.RESET_TRAFFIC,
    responseSchema: ResetNodeTrafficCommand.ResponseSchema,
    routeParamsSchema: ResetNodeTrafficCommand.RequestSchema,
    requestMethod: ResetNodeTrafficCommand.endpointDetails.REQUEST_METHOD,
    rMutationParams: {
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: 'Connection traffic reset successfully',
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: `Reset Connection Traffic`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useBulkSubscriptionConnectionsProfileModification = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.BULK_PROFILE_MODIFICATION,
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
                title: `Bulk Connections Profile Modification`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})

export const useBulkSubscriptionConnectionsActions = createMutationHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.BULK_ACTIONS,
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
                title: `Bulk Connections Actions`,
                message:
                    error instanceof Error ? error.message : `Request failed with unknown error.`,
                color: 'red'
            })
        }
    }
})
