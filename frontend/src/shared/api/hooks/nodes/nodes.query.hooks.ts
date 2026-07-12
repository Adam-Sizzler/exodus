import { z } from 'zod'
import { createQueryKeys } from '@lukemorales/query-key-factory'
import {
    GetAllNodesCommand,
    GetAllNodesTagsCommand,
    GetNodeMetadataCommand,
    GetOneNodeCommand,
    GetPubKeyCommand
} from '@exodus/backend-contract'
import { keepPreviousData } from '@tanstack/react-query'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const nodesQueryKeys = createQueryKeys('nodes', {
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
    },
    getNodeMetadata: (route: GetNodeMetadataCommand.RequestParams) => ({
        queryKey: [route]
    })
})

export const useGetNodes = createGetQueryHook({
    endpoint: GetAllNodesCommand.TSQ_url,
    responseSchema: GetAllNodesCommand.ResponseSchema,
    getQueryKey: () => nodesQueryKeys.getAllNodes.queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get All Nodes')
})

export const useGetNode = createGetQueryHook({
    endpoint: GetOneNodeCommand.TSQ_url,
    responseSchema: GetOneNodeCommand.ResponseSchema.extend({
        response: GetOneNodeCommand.ResponseSchema.shape.response.extend({
            apiSchema: z.enum(['mtls', 'tls']).optional(),
            apiPath: z.string().optional(),
            grpcAuthToken: z.string().optional()
        })
    }),
    routeParamsSchema: GetOneNodeCommand.RequestSchema,
    getQueryKey: ({ route }) => nodesQueryKeys.getNode(route!).queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(5),
        refetchInterval: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get Node')
})
export const useGetPubKey = createGetQueryHook({
    endpoint: GetPubKeyCommand.TSQ_url,
    responseSchema: GetPubKeyCommand.ResponseSchema,
    getQueryKey: () => nodesQueryKeys.getPubKey.queryKey,
    rQueryParams: {
        placeholderData: keepPreviousData,
        refetchOnMount: true,
        staleTime: sToMs(5),
        refetchInterval: sToMs(5)
    },

    errorHandler: (error) => errorHandler(error, 'Get PubKey')
})

export const useGetNodesTags = createGetQueryHook({
    endpoint: GetAllNodesTagsCommand.TSQ_url,
    responseSchema: GetAllNodesTagsCommand.ResponseSchema,
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
