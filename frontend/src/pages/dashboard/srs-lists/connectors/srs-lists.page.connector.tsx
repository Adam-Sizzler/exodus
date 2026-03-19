import { useGetSRSLists } from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui'

import { SRSListsPageComponent } from '../components/srs-lists.page.component'

export function SRSListsPageConnector() {
    const { data, isLoading } = useGetSRSLists()

    if (isLoading || !data) {
        return <LoadingScreen />
    }

    return <SRSListsPageComponent srsLists={data.srsLists} />
}
