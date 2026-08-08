import {
    useGetExternalSquads,
    useGetSubpageConfigs,
    useGetSubscriptionTemplates
} from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui'

import { ExternalSquadsPageComponent } from '../components/external-squads.page.component'

export function ExternalSquadsPageConnector() {
    const { data: externalSquads, isLoading: isExternalSquadsLoading } = useGetExternalSquads()
    const { isLoading: isTemplatesLoading } = useGetSubscriptionTemplates()
    const { isLoading: isSubpageConfigsLoading } = useGetSubpageConfigs()

    if (
        isExternalSquadsLoading ||
        !externalSquads ||
        isTemplatesLoading ||
        isSubpageConfigsLoading
    ) {
        return <LoadingScreen />
    }
    return <ExternalSquadsPageComponent externalSquads={externalSquads.externalSquads} />
}
