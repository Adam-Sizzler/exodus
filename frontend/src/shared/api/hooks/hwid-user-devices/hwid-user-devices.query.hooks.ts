import { createQueryKeys } from '@lukemorales/query-key-factory'
import {
    GetHwidDevicesCommand,
    GetHwidDevicesStatsCommand,
    GetTopUsersByHwidDevicesCommand,
    GetUserHwidDevicesCommand
} from '@exodus/backend-contract'
import { keepPreviousData } from '@tanstack/react-query'

import { sToMs } from '@shared/utils/time-utils'

import { createGetQueryHook, errorHandler } from '../../tsq-helpers'

export const hwidUserDevicesQueryKeys = createQueryKeys('hwid-user-devices', {
    getUserHwidDevices: (route: GetUserHwidDevicesCommand.RequestParam) => ({
        queryKey: [route]
    }),
    getAllHwidDevices: (filters: GetHwidDevicesCommand.RequestQuery) => ({
        queryKey: [filters]
    }),
    getTopUsersByHwidDevices: (filters: GetTopUsersByHwidDevicesCommand.RequestQuery) => ({
        queryKey: [filters]
    }),
    getHwidDevicesStats: {
        queryKey: null
    }
})

const STALE_TIME = 15_000
const REFETCH_INTERVAL = 15_100

export const useGetUserHwidDevices = createGetQueryHook({
    endpoint: GetUserHwidDevicesCommand.TSQ_url,
    responseSchema: GetUserHwidDevicesCommand.ResponseSchema,
    routeParamsSchema: GetUserHwidDevicesCommand.RequestParamSchema,
    getQueryKey: ({ route }) => hwidUserDevicesQueryKeys.getUserHwidDevices(route!).queryKey,
    rQueryParams: {
        refetchInterval: sToMs(60)
    },
    errorHandler: (error) => errorHandler(error, 'Get User HWIDs and Devices')
})

export const useGetHwidDevices = createGetQueryHook({
    endpoint: GetHwidDevicesCommand.TSQ_url,
    responseSchema: GetHwidDevicesCommand.ResponseSchema,
    requestQuerySchema: GetHwidDevicesCommand.RequestQuerySchema,
    getQueryKey: ({ query }) => hwidUserDevicesQueryKeys.getAllHwidDevices(query!).queryKey,
    rQueryParams: {
        staleTime: sToMs(20),
        refetchInterval: sToMs(25),
        placeholderData: keepPreviousData,
        refetchOnMount: true
    },
    errorHandler: (error) => errorHandler(error, 'Get All HWIDs')
})

export const useGetHwidDevicesStats = createGetQueryHook({
    endpoint: GetHwidDevicesStatsCommand.TSQ_url,
    responseSchema: GetHwidDevicesStatsCommand.ResponseSchema,
    getQueryKey: () => hwidUserDevicesQueryKeys.getHwidDevicesStats.queryKey,
    rQueryParams: {
        placeholderData: keepPreviousData,
        staleTime: STALE_TIME,
        refetchInterval: REFETCH_INTERVAL
    },
    errorHandler: (error) => errorHandler(error, 'Get HWIDs Devices Stats')
})

export const useGetTopUsersByHwidDevices = createGetQueryHook({
    endpoint: GetTopUsersByHwidDevicesCommand.TSQ_url,
    responseSchema: GetTopUsersByHwidDevicesCommand.ResponseSchema,
    requestQuerySchema: GetTopUsersByHwidDevicesCommand.RequestQuerySchema,
    getQueryKey: ({ query }) => hwidUserDevicesQueryKeys.getTopUsersByHwidDevices(query!).queryKey,
    rQueryParams: {
        placeholderData: keepPreviousData,
        staleTime: STALE_TIME,
        refetchInterval: REFETCH_INTERVAL
    },
    errorHandler: (error) => errorHandler(error, 'Get Top Users by HWIDs Devices')
})
