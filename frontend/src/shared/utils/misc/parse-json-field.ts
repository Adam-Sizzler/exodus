import consola from 'consola/browser'

export const parseJsonField = (value: unknown): null | object | string => {
    if (value === null || value === undefined || value === '') return null
    if (typeof value === 'object') return value

    try {
        return JSON.parse(value as string)
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

export const parseYamlField = (value: unknown): null | string => {
    if (typeof value !== 'string') return null

    const trimmed = value.trim()
    return trimmed === '' ? null : trimmed
}

export const stringifyYamlField = (value: unknown): string => {
    if (typeof value === 'string') return value

    return ''
}
