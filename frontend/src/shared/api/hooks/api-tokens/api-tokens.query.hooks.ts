import { createQueryKeys } from '@lukemorales/query-key-factory'
import { FindAllApiTokensCommand, GetApiTokenScopesCommand } from '@exodus/backend-contract'
import { keepPreviousData } from '@tanstack/react-query'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const apiTokensQueryKeys = createQueryKeys('apiTokens', {
    getAllApiTokens: {
        queryKey: null
    },
    getApiTokenScopes: {
        queryKey: null
    }
})

export const useGetApiTokens = createGetQueryHook({
    endpoint: FindAllApiTokensCommand.TSQ_url,
    responseSchema: FindAllApiTokensCommand.ResponseSchema,
    getQueryKey: () => apiTokensQueryKeys.getAllApiTokens.queryKey,
    rQueryParams: {
        placeholderData: keepPreviousData,
        refetchOnMount: true
    },
    errorHandler: (error) => errorHandler(error, 'Get All Api Tokens')
})

export const useGetScopes = createGetQueryHook({
    endpoint: GetApiTokenScopesCommand.TSQ_url,
    responseSchema: GetApiTokenScopesCommand.ResponseSchema,
    getQueryKey: () => apiTokensQueryKeys.getApiTokenScopes.queryKey,
    rQueryParams: {},
    errorHandler: (error) => errorHandler(error, 'Get Api Token Scopes')
})
