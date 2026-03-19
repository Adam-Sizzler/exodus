import { createQueryKeys } from '@lukemorales/query-key-factory'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'
import { GetSRSListsResponseSchema } from './srs-lists.schemas'

export const srsListsQueryKeys = createQueryKeys('srsLists', {
    getSRSLists: {
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
