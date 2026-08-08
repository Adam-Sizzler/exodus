import { useGetHttpStats } from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui'
import { sToMs } from '@shared/utils/time-utils'

import HttpStatsPageComponent from '../components/http-stats.page.component'

export function HttpStatsPageConnector() {
    const {
        data: httpStats,
        isLoading,
        refetch,
        isFetching
    } = useGetHttpStats({
        rQueryParams: {
            refetchInterval: sToMs(5)
        }
    })

    if (isLoading || !httpStats) {
        return <LoadingScreen />
    }

    return (
        <HttpStatsPageComponent refetch={refetch} isFetching={isFetching} httpStats={httpStats} />
    )
}
