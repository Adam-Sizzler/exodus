import { createQueryKeys } from '@lukemorales/query-key-factory'
import {
    GetUsersTagsCommand,
    GetUsersCommand,
    GetConnectionKeysByUserIdCommand,
    GetUserAccessibleNodesCommand,
    GetUserByIdCommand,
    GetUserMetadataCommand,
    GetUserSubscriptionRequestHistoryCommand,
    GetRawSubscriptionByShortUuidCommand
} from '@exodus/backend-contract'
import { keepPreviousData } from '@tanstack/react-query'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const usersQueryKeys = createQueryKeys('users', {
    getAllUsers: (filters: GetUsersCommand.RequestQuery) => ({
        queryKey: [filters]
    }),
    getUserById: (route: GetUserByIdCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getConnectionKeysByUserId: (route: GetConnectionKeysByUserIdCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getRawSubscription: (route: GetRawSubscriptionByShortUuidCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getUserTags: {
        queryKey: null
    },
    getUserAccessibleNodes: (route: GetUserAccessibleNodesCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getUserSubscriptionRequestHistory: (
        route: GetUserSubscriptionRequestHistoryCommand.RequestParam
    ) => ({
        queryKey: [route]
    }),
    getUserMetadata: (route: GetUserMetadataCommand.RequestParams) => ({
        queryKey: [route]
    })
})

export const useGetUserById = createGetQueryHook({
    endpoint: GetUserByIdCommand.TSQ_url,
    responseSchema: GetUserByIdCommand.ResponseSchema,
    routeParamsSchema: GetUserByIdCommand.RequestParamSchema,
    getQueryKey: ({ route }) => usersQueryKeys.getUserById(route!).queryKey,
    rQueryParams: {
        staleTime: sToMs(3),
        refetchInterval: sToMs(35)
    },
    errorHandler: (error) => errorHandler(error, 'Get User By User ID')
})

export const useGetUsers = createGetQueryHook({
    endpoint: GetUsersCommand.TSQ_url,
    responseSchema: GetUsersCommand.ResponseSchema,
    requestQuerySchema: GetUsersCommand.RequestQuerySchema,
    getQueryKey: ({ query }) => usersQueryKeys.getAllUsers(query!).queryKey,
    rQueryParams: {
        staleTime: sToMs(20),
        refetchInterval: sToMs(25),
        placeholderData: keepPreviousData,
        refetchOnMount: true
    },
    errorHandler: (error) => errorHandler(error, 'Get All Users')
})

export const useGetUserTags = createGetQueryHook({
    endpoint: GetUsersTagsCommand.TSQ_url,
    responseSchema: GetUsersTagsCommand.ResponseSchema,
    getQueryKey: () => usersQueryKeys.getUserTags.queryKey,
    rQueryParams: {
        staleTime: sToMs(15),
        refetchInterval: sToMs(15)
    },
    errorHandler: (error) => errorHandler(error, 'Get User Tags')
})

export const useGetUserAccessibleNodes = createGetQueryHook({
    endpoint: GetUserAccessibleNodesCommand.TSQ_url,
    responseSchema: GetUserAccessibleNodesCommand.ResponseSchema,
    routeParamsSchema: GetUserAccessibleNodesCommand.RequestParamSchema,
    getQueryKey: ({ route }) => usersQueryKeys.getUserAccessibleNodes(route!).queryKey,
    rQueryParams: {
        staleTime: sToMs(60)
    },
    errorHandler: (error) => errorHandler(error, 'Get User Accessible Nodes')
})

export const useGetUserSubscriptionRequestHistory = createGetQueryHook({
    endpoint: GetUserSubscriptionRequestHistoryCommand.TSQ_url,
    responseSchema: GetUserSubscriptionRequestHistoryCommand.ResponseSchema,
    routeParamsSchema: GetUserSubscriptionRequestHistoryCommand.RequestParamSchema,
    getQueryKey: ({ route }) => usersQueryKeys.getUserSubscriptionRequestHistory(route!).queryKey,
    rQueryParams: {
        staleTime: sToMs(60)
    },
    errorHandler: (error) => errorHandler(error, 'Get User Subscription Request History')
})

export const useGetConnectionKeysByUserId = createGetQueryHook({
    endpoint: GetConnectionKeysByUserIdCommand.TSQ_url,
    responseSchema: GetConnectionKeysByUserIdCommand.ResponseSchema,
    routeParamsSchema: GetConnectionKeysByUserIdCommand.RequestParamSchema,
    getQueryKey: ({ route }) => usersQueryKeys.getConnectionKeysByUserId(route!).queryKey,
    rQueryParams: {
        staleTime: sToMs(4)
    },
    errorHandler: (error) => errorHandler(error, 'Get Connection Keys By User ID')
})

export const useGetUserMetadata = createGetQueryHook({
    endpoint: GetUserMetadataCommand.TSQ_url,
    responseSchema: GetUserMetadataCommand.ResponseSchema,
    routeParamsSchema: GetUserMetadataCommand.RequestParamsSchema,
    getQueryKey: ({ route }) => usersQueryKeys.getUserMetadata(route!).queryKey,
    rQueryParams: {
        staleTime: sToMs(60)
    }
})

export const useGetRawSubscription = createGetQueryHook({
    endpoint: GetRawSubscriptionByShortUuidCommand.TSQ_url,
    responseSchema: GetRawSubscriptionByShortUuidCommand.ResponseSchema,
    routeParamsSchema: GetRawSubscriptionByShortUuidCommand.RequestParamSchema,
    getQueryKey: ({ route }) => usersQueryKeys.getRawSubscription(route!).queryKey,
    rQueryParams: {
        staleTime: sToMs(60)
    }
})
