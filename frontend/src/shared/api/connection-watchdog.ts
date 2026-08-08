import { onlineManager } from '@tanstack/react-query'
import { isAxiosError } from 'axios'

import { instance } from './axios'

const FAILURE_THRESHOLD = 3

let consecutiveFailures = 0

function isNetworkError(error: unknown): boolean {
    return isAxiosError(error) && !error.response && error.code !== 'ERR_CANCELED'
}

export function initConnectionWatchdog() {
    instance.interceptors.response.use(
        (response) => {
            consecutiveFailures = 0
            return response
        },
        (error) => {
            if (onlineManager.isOnline() && isNetworkError(error)) {
                consecutiveFailures += 1
                if (consecutiveFailures >= FAILURE_THRESHOLD) {
                    consecutiveFailures = 0
                    onlineManager.setOnline(false)
                }
            }
            return Promise.reject(error)
        }
    )
}
