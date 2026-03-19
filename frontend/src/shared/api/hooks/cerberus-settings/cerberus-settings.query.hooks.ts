import { GetCerberusSettingsCommand } from '@cerberus/backend-contract'
import { createQueryKeys } from '@lukemorales/query-key-factory'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const cerberusSettingsQueryKeys = createQueryKeys('cerberusSettings', {
    getCerberusSettings: {
        queryKey: null
    }
})

export const useGetCerberusSettings = createGetQueryHook({
    endpoint: GetCerberusSettingsCommand.TSQ_url,
    responseSchema: GetCerberusSettingsCommand.ResponseSchema,
    getQueryKey: () => cerberusSettingsQueryKeys.getCerberusSettings.queryKey,
    rQueryParams: {
        refetchOnMount: false,
        staleTime: sToMs(30)
    },
    errorHandler: (error) => errorHandler(error, 'Get Cerberus Settings')
})
