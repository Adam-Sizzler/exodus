import { useMediaQuery } from '@mantine/hooks'

import { useIsMobile } from '@shared/hooks'
import { HeaderControls } from '@shared/ui/header-buttons'

import { useIsLoadingExodusUpdates, useExodusInfo } from '@entities/dashboard/updates-store'
import { LAYOUT_STYLE, useLayoutStyle } from '@entities/dashboard/view-preferences-store'

import { DASHBOARD_LINKS } from './layout-shared'
import { CompactLayout } from './layout-variants/compact.layout'
import { MobileLayout } from './layout-variants/mobile.layout'
import { SidebarLayout } from './layout-variants/sidebar.layout'

export function MainLayout() {
    const layoutStyle = useLayoutStyle()

    const isMobile = useIsMobile()

    const isHiResDesktop = useMediaQuery(`(min-width: 2048px)`, undefined, {
        getInitialValueInEffect: false
    })

    const exodusInfo = useExodusInfo()
    const isLoadingUpdates = useIsLoadingExodusUpdates()

    const headerControls = (
        <HeaderControls
            {...DASHBOARD_LINKS}
            isGithubLoading={isLoadingUpdates}
            stars={exodusInfo.starsCount || undefined}
            withGithub={!isMobile}
            withPrime
            withRecap={!isMobile}
            withSupport={!isMobile}
            withTelegram={!isMobile}
        />
    )

    if (isMobile) {
        return (
            <MobileLayout
                headerControls={headerControls}
                isSocialButtons={isMobile}
                isLoadingUpdates={isLoadingUpdates}
                exodusInfo={exodusInfo}
            />
        )
    }

    if (layoutStyle === LAYOUT_STYLE.SIDEBAR) {
        return <SidebarLayout headerControls={headerControls} />
    }

    return <CompactLayout headerControls={headerControls} isHiResDesktop={isHiResDesktop} />
}
