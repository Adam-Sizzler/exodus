import { useGetApiTokens, useGetExodusSettings, useGetScopes } from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui/loading-screen'

import { ExodusSettingsPageComponent } from '../components'

export const ExodusSettingsConnector = () => {
    const { data: exodusSettings, isLoading: isExodusSettingsLoading } =
        useGetExodusSettings()
    const { data: apiTokensData, isLoading: isApiTokensLoading } = useGetApiTokens()
    const { data: scopes, isLoading: isScopesLoading } = useGetScopes()

    if (
        isExodusSettingsLoading ||
        isApiTokensLoading ||
        isScopesLoading ||
        !exodusSettings ||
        !apiTokensData ||
        !scopes
    ) {
        return <LoadingScreen />
    }

    return (
        <ExodusSettingsPageComponent
            apiTokensData={apiTokensData}
            exodusSettings={exodusSettings}
        />
    )
}
