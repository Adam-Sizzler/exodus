const normalizeBasePath = (rawBasePath: string): string => {
    const value = rawBasePath.trim()

    if (value === '' || value === '/') {
        return ''
    }

    const withoutSlashes = value.replace(/^\/+|\/+$/g, '')
    if (withoutSlashes === '') {
        return ''
    }

    return `/${withoutSlashes}`
}

const runtimeBasePath =
    typeof window !== 'undefined' ? (window.__EXODUS_RUNTIME__?.basePath ?? '') : ''

export const APP_BASE_PATH = normalizeBasePath(runtimeBasePath)
export const APP_BASE_PATH_WITH_TRAILING_SLASH =
    APP_BASE_PATH === '' ? '/' : `${APP_BASE_PATH}/`

export const withBasePath = (path: string): string => {
    if (!path.startsWith('/')) {
        return path
    }

    if (APP_BASE_PATH === '') {
        return path
    }

    if (path === APP_BASE_PATH || path.startsWith(`${APP_BASE_PATH}/`)) {
        return path
    }

    return `${APP_BASE_PATH}${path}`
}
