import { notifications } from '@mantine/notifications'
import { AxiosError } from 'axios'

const BYPASS_ERROR_STATUSES = [401, 403]

export function errorHandler(error: unknown, title: string) {
    let message = error instanceof Error ? error.message : 'Request failed with unknown error.'

    if (error instanceof AxiosError) {
        if (error.response) {
            if (BYPASS_ERROR_STATUSES.includes(error.response.status)) {
                return
            }

            const errorData = error.response.data
            message =
                typeof errorData === 'string'
                    ? errorData
                    : typeof errorData?.message === 'string'
                      ? errorData.message
                      : typeof errorData?.error === 'string'
                        ? errorData.error
                        : typeof errorData?.error?.message === 'string'
                          ? errorData.error.message
                          : message
        }
    }

    notifications.show({
        title,
        message,
        color: 'red'
    })
}
