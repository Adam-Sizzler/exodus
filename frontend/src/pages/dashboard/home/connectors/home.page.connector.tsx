import { useGetBandwidthStats, useGetCerberusHealth, useGetSystemStats } from '@shared/api/hooks'
import { HomePage } from '@pages/dashboard/home/components'
import { LoadingScreen } from '@shared/ui/loading-screen'

export const HomePageConnector = () => {
    const { data: systemInfo } = useGetSystemStats()
    const { data: bandwidthStats } = useGetBandwidthStats()
    const { data: cerberusHealth } = useGetCerberusHealth()

    if (!systemInfo || !bandwidthStats || !cerberusHealth) {
        return <LoadingScreen />
    }

    return (
        <HomePage
            bandwidthStats={bandwidthStats}
            cerberusHealth={cerberusHealth}
            systemInfo={systemInfo}
        />
    )
}
