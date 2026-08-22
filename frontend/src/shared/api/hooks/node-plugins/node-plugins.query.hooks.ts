import { createQueryKeys } from '@lukemorales/query-key-factory'
import {
    GetNodePluginCommand,
    GetNodePluginsCommand,
    GetSharedListCommand,
    GetSharedListsCommand
} from '@exodus/backend-contract'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const nodePluginsQueryKeys = createQueryKeys('nodePlugins', {
    getNodePlugin: (route: GetNodePluginCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getNodePlugins: {
        queryKey: null
    },
    getSharedList: (route: GetSharedListCommand.RequestParam) => ({
        queryKey: [route]
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
    routeParamsSchema: GetSharedListCommand.RequestParamSchema,
    responseSchema: GetSharedListCommand.ResponseSchema,
    getQueryKey: ({ route }) => nodePluginsQueryKeys.getSharedList(route!).queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get Shared List')
})
