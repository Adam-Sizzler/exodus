import { createQueryKeys } from '@lukemorales/query-key-factory'
import {
    GetSubscriptionTemplateCommand,
    GetSubscriptionTemplatesCommand
} from '@exodus/backend-contract'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const subscriptionTemplateQueryKeys = createQueryKeys('subscriptionTemplate', {
    getSubscriptionTemplate: (route: GetSubscriptionTemplateCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getSubscriptionTemplates: {
        queryKey: null
    }
})

export const useGetSubscriptionTemplate = createGetQueryHook({
    endpoint: GetSubscriptionTemplateCommand.TSQ_url,
    routeParamsSchema: GetSubscriptionTemplateCommand.RequestParamSchema,
    responseSchema: GetSubscriptionTemplateCommand.ResponseSchema,
    getQueryKey: ({ route }) =>
        subscriptionTemplateQueryKeys.getSubscriptionTemplate(route!).queryKey,
    rQueryParams: {
        refetchOnMount: false,
        staleTime: sToMs(30)
    },
    errorHandler: (error) => errorHandler(error, 'Get Subscription Template')
})

export const useGetSubscriptionTemplates = createGetQueryHook({
    endpoint: GetSubscriptionTemplatesCommand.TSQ_url,
    responseSchema: GetSubscriptionTemplatesCommand.ResponseSchema,
    getQueryKey: () => subscriptionTemplateQueryKeys.getSubscriptionTemplates.queryKey,
    rQueryParams: {
        refetchOnMount: false,
        staleTime: sToMs(30)
    },
    errorHandler: (error) => errorHandler(error, 'Get Subscription Templates')
})
