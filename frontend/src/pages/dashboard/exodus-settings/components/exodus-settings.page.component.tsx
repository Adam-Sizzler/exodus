import { FindAllApiTokensCommand, GetExodusSettingsCommand } from '@exodus/backend-contract'
import { useTranslation } from 'node_modules/react-i18next'
import Masonry from 'react-layout-masonry'
import { Container } from '@mantine/core'

import { AuthentificationSettingsCardWidget } from '@widgets/exodus-settings/authentification-settings-card/authentification-settings-card.widget'
import { BrandingSettingsCardWidget } from '@widgets/exodus-settings/branding-settings-card/branding-settings-card.widget'
import { ApiTokensCardWidget } from '@widgets/exodus-settings/api-tokens-card/api-tokens-card.widget'
import { LoadingScreen, Logo, Page, PageHeaderShared } from '@shared/ui'

interface IProps {
    apiTokensData: FindAllApiTokensCommand.Response['response']
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
        !exodusSettings.tgAuthSettings ||
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
                        tgAuthSettings={exodusSettings.tgAuthSettings}
                    />

                    <ApiTokensCardWidget apiTokensData={apiTokensData} />
                    <BrandingSettingsCardWidget
                        brandingSettings={exodusSettings.brandingSettings}
                    />
                </Masonry>
            </Container>
        </Page>
    )
}
