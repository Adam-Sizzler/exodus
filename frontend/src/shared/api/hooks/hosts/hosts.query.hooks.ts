import { createQueryKeys } from '@lukemorales/query-key-factory'
import { GetHostsCommand, GetHostsTagsCommand } from '@exodus/backend-contract'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const hostsQueryKeys = createQueryKeys('hosts', {
    getAllHosts: {
        queryKey: null
    },
    getAllTags: {
        queryKey: null
    }
})

export const useGetHosts = createGetQueryHook({
    endpoint: GetHostsCommand.TSQ_url,
    responseSchema: GetHostsCommand.ResponseSchema,
    getQueryKey: () => hostsQueryKeys.getAllHosts.queryKey,
    rQueryParams: {
        refetchOnMount: true
    },
    errorHandler: (error) => errorHandler(error, 'Get All Hosts')
})

export const useGetHostTags = createGetQueryHook({
    endpoint: GetHostsTagsCommand.TSQ_url,
    responseSchema: GetHostsTagsCommand.ResponseSchema,
    getQueryKey: () => hostsQueryKeys.getAllTags.queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(60),
        refetchInterval: sToMs(60),
        placeholderData: false
    },
    errorHandler: (error) => errorHandler(error, 'Get All Host Tags')
})
