import { Container } from '@mantine/core'
import { GetApiTokensCommand, GetExodusSettingsCommand } from '@exodus/backend-contract'
import { ApiTokensCardWidget } from '@widgets/exodus-settings/api-tokens-card/api-tokens-card.widget'
import { AuthentificationSettingsCardWidget } from '@widgets/exodus-settings/authentification-settings-card/authentification-settings-card.widget'
import { BackendToolsCardWidget } from '@widgets/exodus-settings/backend-tools-card/backend-tools-card.widget'
import { BrandingSettingsCardWidget } from '@widgets/exodus-settings/branding-settings-card/branding-settings-card.widget'
import { VisualSettingsCardWidget } from '@widgets/exodus-settings/visual-settings-card/visual-settings-card.widget'
import { useTranslation } from 'react-i18next'
import Masonry from 'react-layout-masonry'

import { LoadingScreen, Logo, Page, PageHeaderShared } from '@shared/ui'

interface IProps {
    apiTokensData: GetApiTokensCommand.Response['response']
    exodusSettings: GetExodusSettingsCommand.Response['response']
}

export const ExodusSettingsPageComponent = (props: IProps) => {
    const { exodusSettings, apiTokensData } = props

    const { t } = useTranslation()

    if (!exodusSettings || !apiTokensData) {
        return <LoadingScreen />
    }

    if (
        !exodusSettings.oauth2Settings ||
        !exodusSettings.passkeySettings ||
        !exodusSettings.passwordSettings ||
        !exodusSettings.brandingSettings
    ) {
        return <LoadingScreen />
    }

    return (
        <Page title={t('constants.exodus-settings')}>
            <PageHeaderShared icon={<Logo size={24} />} title={t('constants.exodus-settings')} />
            <Container fluid p={0} size="xl">
                <Masonry columns={{ 300: 1, 1400: 2, 2000: 3, 3000: 4 }} gap={16}>
                    <AuthentificationSettingsCardWidget
                        oauth2Settings={exodusSettings.oauth2Settings}
                        passkeySettings={exodusSettings.passkeySettings}
                        passwordSettings={exodusSettings.passwordSettings}
                    />
                    <ApiTokensCardWidget apiTokensData={apiTokensData} />
                    <VisualSettingsCardWidget />
                    <BackendToolsCardWidget />
                    <BrandingSettingsCardWidget
                        brandingSettings={exodusSettings.brandingSettings}
                    />
                </Masonry>
            </Container>
        </Page>
    )
}
