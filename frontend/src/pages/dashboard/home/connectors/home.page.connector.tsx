import { HomePage } from '@pages/dashboard/home/components'

import { useGetBandwidthStats, useGetExodusHealth, useGetSystemStats } from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui/loading-screen'

export const HomePageConnector = () => {
    const { data: systemInfo } = useGetSystemStats()
    const { data: bandwidthStats } = useGetBandwidthStats()
    const { data: exodusHealth } = useGetExodusHealth()

    if (!systemInfo || !bandwidthStats || !exodusHealth) {
        return <LoadingScreen />
    }

    return (
        <HomePage
            bandwidthStats={bandwidthStats}
            exodusHealth={exodusHealth}
            systemInfo={systemInfo}
        />
    )
}
