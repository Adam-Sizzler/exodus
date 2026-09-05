import { createQueryKeys } from '@lukemorales/query-key-factory'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'
import { GetSRSListsResponseSchema, GetSRSListsTagsResponseSchema } from './srs-lists.schemas'

export const srsListsQueryKeys = createQueryKeys('srsLists', {
    getSRSLists: {
        queryKey: null
    },
    getSRSListsTags: {
        queryKey: null
    }
})

export const useGetSRSLists = createGetQueryHook({
    endpoint: '/api/srs-lists',
    responseSchema: GetSRSListsResponseSchema,
    getQueryKey: () => srsListsQueryKeys.getSRSLists.queryKey,
    rQueryParams: {
        refetchOnMount: true,
        refetchInterval: sToMs(300),
        staleTime: sToMs(30)
    },
    errorHandler: (error) => errorHandler(error, 'Get SRS Lists')
})

export const useGetSRSListsTags = createGetQueryHook({
    endpoint: '/api/srs-lists/tags',
    responseSchema: GetSRSListsTagsResponseSchema,
    getQueryKey: () => srsListsQueryKeys.getSRSListsTags.queryKey,
    rQueryParams: {
        refetchOnMount: true,
        staleTime: sToMs(30)
    },
    errorHandler: (error) => errorHandler(error, 'Get SRS Lists Tags')
})
