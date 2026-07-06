import { modals } from '@mantine/modals'
import ColorHash from 'color-hash'
import { useTranslation } from 'react-i18next'
import { PiAppWindowDuotone } from 'react-icons/pi'

import { TopLeaderboardCardShared } from '@shared/ui/leaderboard-item-card'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { IPlatformDatum } from './hwid-inspector-metrics.types'
import { formatCount } from './hwid-inspector-metrics.utils'

function PlatformAppsModalContent({ platform }: { platform: IPlatformDatum }) {
    const { t } = useTranslation()
    const ch = new ColorHash({ lightness: [0.65, 0.65, 0.65] })

    return (
        <TopLeaderboardCardShared
            emptyText={t('hwid-inspector-metrics.widget.no-app-data')}
            formatValue={formatCount}
            isLoading={false}
            items={platform.allApps.map((app) => ({
                color: ch.hex(app.app),
                name: app.app,
                total: app.count
            }))}
            maxHeight={420}
            wrapper={(children) => children}
        />
    )
}

export function openPlatformAppsModal(platform: IPlatformDatum, isMobile: boolean) {
    modals.open({
        centered: true,
        size: 'lg',
        fullScreen: isMobile,
        title: (
            <BaseOverlayHeader
                IconComponent={PiAppWindowDuotone}
                iconColor="cyan"
                iconVariant="soft"
                subtitle={formatCount(platform.value)}
                title={platform.name}
            />
        ),
        children: <PlatformAppsModalContent platform={platform} />
    })
}
