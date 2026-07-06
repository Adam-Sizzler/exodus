import { createQueryKeys } from '@lukemorales/query-key-factory'
import {
    GetPasskeyAuthenticationOptionsCommand,
    GetStatusCommand
} from '@exodus/backend-contract'
import { keepPreviousData } from '@tanstack/react-query'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const authQueryKeys = createQueryKeys('auth', {
    getAuthStatus: {
        queryKey: null
    },
    getPasskeyAuthenticationOptions: {
        queryKey: null
    }
})

export const useGetAuthStatus = createGetQueryHook({
    endpoint: GetStatusCommand.TSQ_url,
    responseSchema: GetStatusCommand.ResponseSchema,
    getQueryKey: () => authQueryKeys.getAuthStatus.queryKey,
    rQueryParams: {
        refetchOnMount: false,
        placeholderData: keepPreviousData
    },
    errorHandler: (error) => errorHandler(error, 'Authentication Error')
})

export const usePasskeyAuthenticationOptions = createGetQueryHook({
    endpoint: GetPasskeyAuthenticationOptionsCommand.TSQ_url,
    responseSchema: GetPasskeyAuthenticationOptionsCommand.ResponseSchema,
    getQueryKey: () => authQueryKeys.getPasskeyAuthenticationOptions.queryKey,
    rQueryParams: {
        enabled: false
    },
    errorHandler: (error) => errorHandler(error, 'Authentication Error')
})
