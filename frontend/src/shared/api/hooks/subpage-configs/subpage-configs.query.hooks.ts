import { createQueryKeys } from '@lukemorales/query-key-factory'
import { GetSubpageConfigCommand, GetSubpageConfigsCommand } from '@exodus/backend-contract'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const subpageConfigsQueryKeys = createQueryKeys('subpageConfigs', {
    getSubpageConfig: (route: GetSubpageConfigCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getSubpageConfigs: {
        queryKey: null
    }
})

export const useGetSubscriptionPageConfig = createGetQueryHook({
    endpoint: GetSubpageConfigCommand.TSQ_url,
    routeParamsSchema: GetSubpageConfigCommand.RequestParamSchema,
    responseSchema: GetSubpageConfigCommand.ResponseSchema,
    getQueryKey: ({ route }) => subpageConfigsQueryKeys.getSubpageConfig(route!).queryKey,
    rQueryParams: {
        refetchOnMount: false,
        staleTime: sToMs(5)
    },
    errorHandler: (error) => errorHandler(error, 'Get Subscription Page Config')
})

export const useGetSubpageConfigs = createGetQueryHook({
    endpoint: GetSubpageConfigsCommand.TSQ_url,
    responseSchema: GetSubpageConfigsCommand.ResponseSchema,
    getQueryKey: () => subpageConfigsQueryKeys.getSubpageConfigs.queryKey,
    rQueryParams: {
        refetchOnMount: false,
        staleTime: sToMs(15)
    },
    errorHandler: (error) => errorHandler(error, 'Get Subscription Page Configs')
})
