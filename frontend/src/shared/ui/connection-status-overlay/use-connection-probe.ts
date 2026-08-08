import { GetExodusHealthCommand, GetStatusCommand } from '@exodus/backend-contract'
import { onlineManager, useQuery } from '@tanstack/react-query'
import { isAxiosError } from 'axios'
import { useSyncExternalStore } from 'react'

import { hasAuthorizationToken, instance } from '@shared/api'

const INITIAL_RETRY_DELAY_MS = 1_000
const MAX_RETRY_DELAY_MS = 30_000
const PROBE_TIMEOUT_MS = 5_000

function subscribeToOnlineManager(onStoreChange: () => void) {
    return onlineManager.subscribe(onStoreChange)
}

function getOnlineSnapshot() {
    return onlineManager.isOnline()
}

export function useConnectionProbe() {
    const isOnline = useSyncExternalStore(subscribeToOnlineManager, getOnlineSnapshot)

    useQuery({
        queryKey: ['connection-probe'],
        enabled: !isOnline,
        networkMode: 'always',
        retry: true,
        retryDelay: (attemptIndex) =>
            Math.min(INITIAL_RETRY_DELAY_MS * 2 ** attemptIndex, MAX_RETRY_DELAY_MS),
        staleTime: 0,
        gcTime: 0,
        queryFn: async ({ signal }) => {
            const probeUrl = hasAuthorizationToken()
                ? GetExodusHealthCommand.TSQ_url
                : GetStatusCommand.TSQ_url

            try {
                await instance.get(probeUrl, { signal, timeout: PROBE_TIMEOUT_MS })
            } catch (error) {
                if (!isAxiosError(error) || !error.response) {
                    throw error
                }
            }

            onlineManager.setOnline(true)
            return true
        }
    })

    return { isOnline }
}
