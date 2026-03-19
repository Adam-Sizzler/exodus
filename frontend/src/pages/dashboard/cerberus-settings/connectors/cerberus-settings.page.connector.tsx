import { useGetApiTokens, useGetCerberusSettings } from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui/loading-screen'

import { CerberusSettingsPageComponent } from '../components'

export const CerberusSettingsConnector = () => {
    const { data: cerberusSettings, isLoading: isCerberusSettingsLoading } =
        useGetCerberusSettings()
    const { data: apiTokensData, isLoading: isApiTokensLoading } = useGetApiTokens()

    if (isCerberusSettingsLoading || isApiTokensLoading || !cerberusSettings || !apiTokensData) {
        return <LoadingScreen />
    }

    return (
        <CerberusSettingsPageComponent
            apiTokensData={apiTokensData}
            cerberusSettings={cerberusSettings}
        />
    )
}
