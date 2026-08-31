import { createQueryKeys } from '@lukemorales/query-key-factory'
import {
    GetNodePluginCommand,
    GetNodePluginsCommand,
    GetNodePluginsTagsCommand,
    GetSharedListCommand,
    GetSharedListsCommand
} from '@exodus/backend-contract'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const nodePluginsQueryKeys = createQueryKeys('nodePlugins', {
    getNodePluginsTags: {
        queryKey: null
    },
    getNodePlugin: (route: GetNodePluginCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getNodePlugins: {
        queryKey: null
    },
    getSharedList: (query: GetSharedListCommand.RequestQuery) => ({
        queryKey: [query]
    }),
    getSharedLists: {
        queryKey: null
    }
})

export const useGetNodePlugin = createGetQueryHook({
    endpoint: GetNodePluginCommand.TSQ_url,
    routeParamsSchema: GetNodePluginCommand.RequestParamSchema,
    responseSchema: GetNodePluginCommand.ResponseSchema,
    getQueryKey: ({ route }) => nodePluginsQueryKeys.getNodePlugin(route!).queryKey,
    rQueryParams: {
        refetchOnMount: false,
        staleTime: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get Node Plugin')
})

export const useGetNodePlugins = createGetQueryHook({
    endpoint: GetNodePluginsCommand.TSQ_url,
    responseSchema: GetNodePluginsCommand.ResponseSchema,
    getQueryKey: () => nodePluginsQueryKeys.getNodePlugins.queryKey,
    rQueryParams: {
        refetchOnMount: false,
        staleTime: sToMs(15)
    },
    errorHandler: (error) => errorHandler(error, 'Get Node Plugins')
})

export const useGetSharedLists = createGetQueryHook({
    endpoint: GetSharedListsCommand.TSQ_url,
    responseSchema: GetSharedListsCommand.ResponseSchema,
    getQueryKey: () => nodePluginsQueryKeys.getSharedLists.queryKey,
    rQueryParams: {
        refetchOnMount: false,
        staleTime: sToMs(15)
    },
    errorHandler: (error) => errorHandler(error, 'Get Shared Lists')
})

export const useGetSharedList = createGetQueryHook({
    endpoint: GetSharedListCommand.TSQ_url,
    requestQuerySchema: GetSharedListCommand.RequestQuerySchema,
    responseSchema: GetSharedListCommand.ResponseSchema,
    getQueryKey: ({ query }) => nodePluginsQueryKeys.getSharedList(query!).queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get Shared List')
})

export const useGetNodePluginsTags = createGetQueryHook({
    endpoint: GetNodePluginsTagsCommand.TSQ_url,
    responseSchema: GetNodePluginsTagsCommand.ResponseSchema,
    getQueryKey: () => nodePluginsQueryKeys.getNodePluginsTags.queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(30)
    },
    errorHandler: (error) => errorHandler(error, 'Get NodePlugins Tags')
})
