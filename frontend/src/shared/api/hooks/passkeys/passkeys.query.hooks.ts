import { createQueryKeys } from '@lukemorales/query-key-factory'
import {
    GetPasskeysCommand,
    GetPasskeyRegistrationOptionsCommand
} from '@exodus/backend-contract'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const passkeysQueryKeys = createQueryKeys('passkeys', {
    getPasskeys: {
        queryKey: null
    },
    getPasskeyRegistrationOptions: {
        queryKey: null
    }
})

export const usePasskeyRegistrationOptions = createGetQueryHook({
    endpoint: GetPasskeyRegistrationOptionsCommand.TSQ_url,
    responseSchema: GetPasskeyRegistrationOptionsCommand.ResponseSchema,
    getQueryKey: () => passkeysQueryKeys.getPasskeyRegistrationOptions.queryKey,
    rQueryParams: {
        enabled: false
    },
    errorHandler: (error) => errorHandler(error, 'Get Passkey Registration Options')
})

export const useGetPasskeys = createGetQueryHook({
    endpoint: GetPasskeysCommand.TSQ_url,
    responseSchema: GetPasskeysCommand.ResponseSchema,
    getQueryKey: () => passkeysQueryKeys.getPasskeys.queryKey,
    rQueryParams: {
        refetchOnMount: false
    },
    errorHandler: (error) => errorHandler(error, 'Get All Passkeys')
})
