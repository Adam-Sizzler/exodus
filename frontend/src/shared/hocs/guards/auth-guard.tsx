import { useLayoutEffect } from 'react'
import { Navigate, Outlet, useLocation } from 'react-router'

import { useGetAuthStatus } from '@shared/api/hooks'
import { ROUTES } from '@shared/constants/routes'
import { useAuth } from '@shared/hooks'
import { LoadingProgress } from '@shared/ui/loading-screen'
import { consumeReturnTo, saveReturnTo } from '@shared/utils/return-to.util'

import { useUpdatesStoreActions } from '@entities/dashboard/updates-store'

export function AuthGuard() {
    const location = useLocation()

    const { isAuthenticated, isInitialized } = useAuth()

    const { isLoading } = useGetAuthStatus()
    const updatesStoreActions = useUpdatesStoreActions()

    useLayoutEffect(() => {
        updatesStoreActions.getExodusInfo()
    }, [])

    if (!isInitialized || isLoading) {
        return <LoadingProgress />
    }

    if (!isAuthenticated) {
        if (location.pathname.includes(ROUTES.AUTH.ROOT)) {
            return <Outlet />
        }
        saveReturnTo(location.pathname + location.search)
        return <Navigate replace to={ROUTES.AUTH.LOGIN} />
    }

    if (isAuthenticated) {
        if (location.pathname.includes(ROUTES.DASHBOARD.ROOT)) {
            return <Outlet />
        }
        return <Navigate replace to={consumeReturnTo() ?? ROUTES.DASHBOARD.ROOT} />
    }

    return <Outlet />
}
