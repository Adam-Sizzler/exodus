import {
    GetAllNodesCommand,
    GetAllNodesTagsCommand,
    GetOneNodeCommand,
    GetPubKeyCommand
} from '@exodus/backend-contract'
import { createQueryKeys } from '@lukemorales/query-key-factory'
import { keepPreviousData } from '@tanstack/react-query'
import { z } from 'zod'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

const nodeSystemInterfaceSchema = z.object({
    interface: z.string(),
    rxBytesPerSec: z.number().optional().default(0),
    txBytesPerSec: z.number().optional().default(0),
    rxTotal: z.number().optional().default(0),
    txTotal: z.number().optional().default(0)
})

const nodeSystemSchema = z.object({
    info: z.object({
        arch: z.string().optional().default('unknown'),
        cpus: z.number().optional().default(1),
        cpuModel: z.string().optional().default('unknown'),
        memoryTotal: z.number().optional().default(0),
        hostname: z.string().optional().default('unknown'),
        platform: z.string().optional().default('unknown'),
        release: z.string().optional().default('unknown'),
        type: z.string().optional().default('unknown'),
        version: z.string().optional().default('unknown'),
        networkInterfaces: z.array(z.string()).optional().default([])
    }),
    stats: z.object({
        memoryFree: z.number().optional().default(0),
        memoryUsed: z.number().optional().default(0),
        uptime: z.number().optional().default(0),
        loadAvg: z.array(z.number()).optional().default([0, 0, 0]),
        interface: nodeSystemInterfaceSchema.nullable().optional().default(null)
    })
})

const nodePluginSchema = z.object({
    uuid: z.string(),
    name: z.string(),
    pluginConfig: z.record(z.unknown()).default({}),
    viewPosition: z.number().optional().default(0)
})

const nodeResponseSchema = GetOneNodeCommand.ResponseSchema.shape.response.extend({
    apiSchema: z.string().optional().default('mtls'),
    apiPath: z.string().optional().default('/'),
    grpcAuthToken: z.string().optional().default(''),
    activePluginUuid: z.string().nullable().optional().default(null),
    system: nodeSystemSchema.nullable().optional().default(null),
    versions: z
        .object({
            singbox: z.string().optional().default('unknown'),
            node: z.string().optional().default('unknown')
        })
        .nullable()
        .optional()
        .default(null)
})

const getAllNodesResponseSchema = z.object({
    response: z.array(nodeResponseSchema)
})

const getOneNodeResponseSchema = z.object({
    response: nodeResponseSchema
})

const getPubKeyResponseSchema = z.object({
    response: z.object({
        pubKey: z.string(),
        grpcToken: z.string().optional().default('')
    })
})

export type NodeKeygenResponse = z.infer<typeof getPubKeyResponseSchema>['response']
export type NodePluginResponse = z.infer<typeof nodePluginSchema>
export type NodeResponse = z.infer<typeof nodeResponseSchema>

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
    }
})

export const useGetNodes = createGetQueryHook({
    endpoint: GetAllNodesCommand.TSQ_url,
    responseSchema: getAllNodesResponseSchema,
    getQueryKey: () => nodesQueryKeys.getAllNodes.queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get All Nodes')
})

export const useGetNode = createGetQueryHook({
    endpoint: GetOneNodeCommand.TSQ_url,
    responseSchema: getOneNodeResponseSchema,
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
    responseSchema: getPubKeyResponseSchema,
    getQueryKey: () => nodesQueryKeys.getPubKey.queryKey,
    rQueryParams: {
        placeholderData: keepPreviousData,
        refetchOnMount: true,
        staleTime: sToMs(5)
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
