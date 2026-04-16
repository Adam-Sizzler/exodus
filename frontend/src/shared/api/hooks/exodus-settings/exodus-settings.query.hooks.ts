import { GetExodusSettingsCommand } from '@exodus/backend-contract'
import { createQueryKeys } from '@lukemorales/query-key-factory'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const exodusSettingsQueryKeys = createQueryKeys('exodusSettings', {
    getExodusSettings: {
        queryKey: null
    }
})

export const useGetExodusSettings = createGetQueryHook({
    endpoint: GetExodusSettingsCommand.TSQ_url,
    responseSchema: GetExodusSettingsCommand.ResponseSchema,
    getQueryKey: () => exodusSettingsQueryKeys.getExodusSettings.queryKey,
    rQueryParams: {
        refetchOnMount: false,
        staleTime: sToMs(30)
    },
    errorHandler: (error) => errorHandler(error, 'Get Exodus Settings')
})
