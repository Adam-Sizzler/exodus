import { ROUTES } from '@shared/constants/routes'

const RETURN_TO_KEY = 'returnTo'

const isSafeReturnPath = (path: null | string): path is string => {
    return Boolean(path) && path!.startsWith(ROUTES.DASHBOARD.ROOT)
}

export const saveReturnTo = (path: string) => {
    if (!isSafeReturnPath(path)) {
        return
    }

    try {
        sessionStorage.setItem(RETURN_TO_KEY, path)
    } catch {
        //
    }
}

export const consumeReturnTo = (): null | string => {
    try {
        const path = sessionStorage.getItem(RETURN_TO_KEY)
        sessionStorage.removeItem(RETURN_TO_KEY)

        return isSafeReturnPath(path) ? path : null
    } catch {
        return null
    }
}
