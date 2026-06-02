import { consola } from 'consola/browser'
import { isAxiosError } from 'axios'
import { ZodError } from 'zod'

/** Handle request errors */
export function handleRequestError(error: unknown) {
    if (isAxiosError(error)) {
        const errorData = error.response?.data
        const message =
            typeof errorData === 'string'
                ? errorData
                : typeof errorData?.message === 'string'
                  ? errorData.message
                  : typeof errorData?.error === 'string'
                    ? errorData.error
                    : typeof errorData?.error?.message === 'string'
                      ? errorData.error.message
                      : 'Request failed'
        const enhancedError = new Error(message)
        enhancedError.cause = errorData
        throw enhancedError
    }

    if (error instanceof ZodError) {
        consola.error(error.format())
    }

    consola.log(error)

    throw error
}
