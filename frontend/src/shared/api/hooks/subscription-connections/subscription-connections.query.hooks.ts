import { GetNodesTagsCommand, GetNodeCommand } from '@exodus/backend-contract'
import { createQueryKeys } from '@lukemorales/query-key-factory'
import { keepPreviousData } from '@tanstack/react-query'
import { z } from 'zod'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

const SUBSCRIPTION_CONNECTIONS_API = {
    GET_ALL: '/api/subscription-connections',
    GET_ONE: '/api/subscription-connections/:uuid',
    GET_TAGS: '/api/subscription-connections/tags'
} as const

const dateTimeSchema = z.string().datetime().transform((str) => new Date(str))

const subscriptionConnectionConfigInboundSchema = z.object({
    uuid: z.string().uuid(),
    profileUuid: z.string().uuid(),
    tag: z.string(),
    type: z.string(),
    network: z.string().nullable(),
    security: z.string().nullable(),
    port: z.number().nullable(),
    rawInbound: z.unknown().nullable()
})

const subscriptionConnectionProviderSchema = z.object({
    uuid: z.string().uuid(),
    name: z.string(),
    faviconLink: z.string().nullable(),
    loginUrl: z.string().nullable(),
    createdAt: dateTimeSchema.optional(),
    updatedAt: dateTimeSchema.optional()
})

const subscriptionConnectionNetworkInterfaceSchema = z.object({
    interface: z.string(),
    rxBytesPerSec: z.number(),
    txBytesPerSec: z.number(),
    rxTotal: z.number(),
    txTotal: z.number()
})

const subscriptionConnectionSystemSchema = z
    .object({
        info: z.object({
            arch: z.string(),
            cpus: z.number().int(),
            cpuModel: z.string(),
            memoryTotal: z.number(),
            hostname: z.string(),
            platform: z.string(),
            release: z.string(),
            type: z.string(),
            version: z.string(),
            networkInterfaces: z.array(z.string())
        }),
        stats: z.object({
            memoryFree: z.number(),
            memoryUsed: z.number(),
            uptime: z.number(),
            loadAvg: z.array(z.number()),
            interface: subscriptionConnectionNetworkInterfaceSchema.nullable()
        })
    })
    .nullable()
    .optional()
    .default(null)

export const subscriptionConnectionNodeSchema = z.object({
    uuid: z.string().uuid(),
    name: z.string(),
    address: z.string(),
    publicDomain: z.string().nullable().optional(),
    port: z.number().int().nullable(),
    apiSchema: z.string().optional().default('mtls'),
    apiPath: z.string().optional().default('/'),
    grpcAuthToken: z.string().optional().default(''),
    subpageConfigUuid: z.string().uuid().nullable().optional(),
    isConnected: z.boolean(),
    isDisabled: z.boolean(),
    isConnecting: z.boolean(),
    lastStatusChange: dateTimeSchema.nullable(),
    lastStatusMessage: z.string().nullable(),
    singboxVersion: z.string().nullable().optional(),
    nodeVersion: z.string().nullable().optional(),
    singboxUptime: z.number(),
    isTrafficTrackingActive: z.boolean(),
    trafficResetDay: z.number().int().nullable(),
    trafficLimitBytes: z.number().nullable(),
    trafficUsedBytes: z.number().nullable(),
    notifyPercent: z.number().int().nullable(),
    usersOnline: z.number().int().nullable(),
    viewPosition: z.number().int(),
    countryCode: z.string(),
    consumptionMultiplier: z.number(),
    tags: z.array(z.string()),
    cpuCount: z.number().int().nullable().optional(),
    cpuModel: z.string().nullable().optional(),
    totalRam: z.string().nullable().optional(),
    system: subscriptionConnectionSystemSchema,
    versions: z
        .object({
            singbox: z.string().optional().default(''),
            node: z.string().optional().default('')
        })
        .nullable()
        .optional()
        .default(null),
    createdAt: dateTimeSchema,
    updatedAt: dateTimeSchema,
    configProfile: z.object({
        activeConfigProfileUuid: z.string().uuid().nullable(),
        activeInbounds: z.array(subscriptionConnectionConfigInboundSchema)
    }),
    providerUuid: z.string().uuid().nullable().optional(),
    provider: subscriptionConnectionProviderSchema.nullable()
})

export const getAllSubscriptionConnectionsResponseSchema = z.object({
    response: z.array(subscriptionConnectionNodeSchema)
})

export const getOneSubscriptionConnectionResponseSchema = z.object({
    response: subscriptionConnectionNodeSchema
})

export const subscriptionConnectionEventSentResponseSchema = z.object({
    response: z.object({
        eventSent: z.boolean()
    })
})

export const deleteSubscriptionConnectionResponseSchema = z.object({
    response: z.object({
        isDeleted: z.boolean()
    })
})

const getSubscriptionConnectionsPubKeyResponseSchema = z.object({
    response: z
        .object({
            secretKey: z.string().optional(),
            pubKey: z.string().optional(),
            grpcToken: z.string().optional().default('')
        })
        .transform((data) => ({
            secretKey: data.secretKey || data.pubKey || '',
            pubKey: data.pubKey || data.secretKey || '',
            grpcToken: data.grpcToken || ''
        }))
})

export type SubscriptionConnectionKeygenResponse = z.infer<
    typeof getSubscriptionConnectionsPubKeyResponseSchema
>['response']
export type SubscriptionConnectionResponse = z.infer<typeof subscriptionConnectionNodeSchema>

export const subscriptionConnectionsQueryKeys = createQueryKeys('subscriptionConnections', {
    getAllNodes: {
        queryKey: null
    },
    getNode: (route: GetNodeCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getPubKey: {
        queryKey: null
    },
    getAllTags: {
        queryKey: null
    }
})

export const useGetSubNodes = createGetQueryHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.GET_ALL,
    responseSchema: getAllSubscriptionConnectionsResponseSchema,
    getQueryKey: () => subscriptionConnectionsQueryKeys.getAllNodes.queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get All Subscription Connections')
})

export const useGetSubscriptionConnection = createGetQueryHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.GET_ONE,
    responseSchema: getOneSubscriptionConnectionResponseSchema,
    routeParamsSchema: GetNodeCommand.RequestParamSchema,
    getQueryKey: ({ route }) => subscriptionConnectionsQueryKeys.getNode(route!).queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(5),
        refetchInterval: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get Subscription Connection')
})

export const useGetSubscriptionConnectionsPubKey = createGetQueryHook({
    endpoint: '/api/keygen',
    responseSchema: getSubscriptionConnectionsPubKeyResponseSchema,
    getQueryKey: () => subscriptionConnectionsQueryKeys.getPubKey.queryKey,
    rQueryParams: {
        placeholderData: keepPreviousData,
        refetchOnMount: true,
        staleTime: sToMs(5)
    },

    errorHandler: (error) => errorHandler(error, 'Get Subscription Connection PubKey')
})

export const useGetSubscriptionConnectionsTags = createGetQueryHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.GET_TAGS,
    responseSchema: GetNodesTagsCommand.ResponseSchema,
    getQueryKey: () => subscriptionConnectionsQueryKeys.getAllTags.queryKey,
    rQueryParams: {
        staleTime: 0
    },
    errorHandler: (error) => errorHandler(error, 'Get All Subscription Connection Tags')
})
