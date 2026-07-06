import { formatInt, formatPercentage } from '@shared/utils/misc'

import {
    APP_MIN_SHARE,
    CHART_COLORS,
    MAX_APPS,
    MAX_PLATFORMS
} from './hwid-inspector-metrics.constants'
import { IPlatformApp, IPlatformDatum } from './hwid-inspector-metrics.types'

export const formatCount = (value: number) => formatInt(value, { thousandSeparator: ' ' })

export const formatShare = (value: number, total: number) =>
    formatPercentage(total > 0 ? (value / total) * 100 : 0, { precision: 1 })

export function buildPlatformData(
    byPlatform: { byApp: IPlatformApp[]; count: number; platform: string }[]
): IPlatformDatum[] {
    return byPlatform
        .filter((platform) => platform.count > 0)
        .slice(0, MAX_PLATFORMS)
        .map((platform, index) => {
            const allApps = platform.byApp.map((app) => ({
                app: app.app || 'Unknown',
                count: app.count
            }))

            const apps = allApps
                .filter(
                    (app, appIndex) => appIndex === 0 || app.count >= platform.count * APP_MIN_SHARE
                )
                .slice(0, MAX_APPS)

            return {
                name: platform.platform || 'Unknown',
                value: platform.count,
                color: CHART_COLORS[index % CHART_COLORS.length],
                apps,
                allApps
            }
        })
}
