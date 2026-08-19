import { z } from 'zod'
import { createQueryKeys } from '@lukemorales/query-key-factory'
import {
    GetNodeMetadataCommand,
    GetNodeCommand,
    GetNodeSecretKeyCommand,
    GetNodesTagsCommand,
    GetNodesCommand
} from '@exodus/backend-contract'
import { keepPreviousData } from '@tanstack/react-query'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const nodesQueryKeys = createQueryKeys('nodes', {
    getAllNodes: {
        queryKey: null
    },
    getNode: (route: GetNodeCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getNodeSecretKey: {
        queryKey: null
    },
    getAllTags: {
        queryKey: null
    },
    getNodeMetadata: (route: GetNodeMetadataCommand.RequestParams) => ({
        queryKey: [route]
    })
})

export const useGetNodes = createGetQueryHook({
    endpoint: GetNodesCommand.TSQ_url,
    responseSchema: GetNodesCommand.ResponseSchema,
    getQueryKey: () => nodesQueryKeys.getAllNodes.queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get All Nodes')
})

export const useGetNode = createGetQueryHook({
    endpoint: GetNodeCommand.TSQ_url,
    responseSchema: GetNodeCommand.ResponseSchema.extend({
        response: GetNodeCommand.ResponseSchema.shape.response.extend({
            apiSchema: z.enum(['mtls', 'tls']).optional(),
            apiPath: z.string().optional(),
            grpcAuthToken: z.string().optional()
        })
    }),
    routeParamsSchema: GetNodeCommand.RequestParamSchema,
    getQueryKey: ({ route }) => nodesQueryKeys.getNode(route!).queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(5),
        refetchInterval: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get Node')
})
export const useGetNodeSecretKey = createGetQueryHook({
    endpoint: GetNodeSecretKeyCommand.TSQ_url,
    responseSchema: GetNodeSecretKeyCommand.ResponseSchema,
    getQueryKey: () => nodesQueryKeys.getNodeSecretKey.queryKey,
    rQueryParams: {
        placeholderData: keepPreviousData,
        refetchOnMount: true,
        staleTime: sToMs(5),
        refetchInterval: sToMs(5)
    },

    errorHandler: (error) => errorHandler(error, 'Get SecretKey')
})

export const useGetNodesTags = createGetQueryHook({
    endpoint: GetNodesTagsCommand.TSQ_url,
    responseSchema: GetNodesTagsCommand.ResponseSchema,
    getQueryKey: () => nodesQueryKeys.getAllTags.queryKey,
    rQueryParams: {
        staleTime: 0
    },
    errorHandler: (error) => errorHandler(error, 'Get All Nodes Tags')
})

export const useGetNodeMetadata = createGetQueryHook({
    endpoint: GetNodeMetadataCommand.TSQ_url,
    responseSchema: GetNodeMetadataCommand.ResponseSchema,
    routeParamsSchema: GetNodeMetadataCommand.RequestParamsSchema,
    getQueryKey: ({ route }) => nodesQueryKeys.getNodeMetadata(route!).queryKey,
    rQueryParams: {
        staleTime: sToMs(60)
    }
})
