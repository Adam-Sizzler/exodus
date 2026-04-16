import { useGetApiTokens, useGetExodusSettings } from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui/loading-screen'

import { ExodusSettingsPageComponent } from '../components'

export const ExodusSettingsConnector = () => {
    const { data: exodusSettings, isLoading: isExodusSettingsLoading } =
        useGetExodusSettings()
    const { data: apiTokensData, isLoading: isApiTokensLoading } = useGetApiTokens()

    if (isExodusSettingsLoading || isApiTokensLoading || !exodusSettings || !apiTokensData) {
        return <LoadingScreen />
    }

    return (
        <ExodusSettingsPageComponent
            apiTokensData={apiTokensData}
            exodusSettings={exodusSettings}
        />
    )
}
