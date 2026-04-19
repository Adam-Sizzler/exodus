import {
    GetAllNodesTagsCommand,
    GetOneNodeCommand
} from '@exodus/backend-contract'
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

const subscriptionConnectionNodeSchema = GetOneNodeCommand.ResponseSchema.shape.response.extend({
    apiSchema: z.string().optional().default('mtls'),
    apiPath: z.string().optional().default('/'),
    publicDomain: z.string().nullable().optional(),
    subpageConfigUuid: z.string().uuid().nullable().optional()
})

const getAllSubscriptionConnectionsResponseSchema = z.object({
    response: z.array(subscriptionConnectionNodeSchema)
})

const getOneSubscriptionConnectionResponseSchema = z.object({
    response: subscriptionConnectionNodeSchema
})

const getSubscriptionConnectionsPubKeyResponseSchema = z.object({
    response: z.object({
        pubKey: z.string(),
        grpcToken: z.string().optional().default('')
    })
})

export const subscriptionConnectionsQueryKeys = createQueryKeys('subscriptionConnections', {
    getAllNodes: {
        queryKey: null
    },
    getNode: (route: GetOneNodeCommand.Request) => ({
        queryKey: [route]
    }),
    getPubKey: {
        queryKey: null
    },
    getAllTags: {
        queryKey: null
    }
})

export const useGetSubscriptionConnections = createGetQueryHook({
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
    routeParamsSchema: GetOneNodeCommand.RequestSchema,
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
        staleTime: sToMs(5),
        refetchInterval: sToMs(5)
    },

    errorHandler: (error) => errorHandler(error, 'Get Subscription Connection PubKey')
})

export const useGetSubscriptionConnectionsTags = createGetQueryHook({
    endpoint: SUBSCRIPTION_CONNECTIONS_API.GET_TAGS,
    responseSchema: GetAllNodesTagsCommand.ResponseSchema,
    getQueryKey: () => subscriptionConnectionsQueryKeys.getAllTags.queryKey,
    rQueryParams: {
        staleTime: 0
    },
    errorHandler: (error) => errorHandler(error, 'Get All Subscription Connection Tags')
})
