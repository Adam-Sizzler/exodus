import consola from 'consola/browser'

export const parseJsonField = (value: unknown): null | object | string => {
    if (value === '') return null

    try {
        return JSON.parse(value as unknown as string)
    } catch (error) {
        consola.error(error)
        return null
    }
}

export const stringifyJsonField = (value: unknown): string => {
    if (typeof value === 'object' && value !== null) {
        return JSON.stringify(value, null, 2)
    }

    return ''
}
