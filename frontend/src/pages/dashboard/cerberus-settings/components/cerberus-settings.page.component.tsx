import { FindAllApiTokensCommand, GetCerberusSettingsCommand } from '@cerberus/backend-contract'
import { useTranslation } from 'node_modules/react-i18next'
import Masonry from 'react-layout-masonry'
import { Container } from '@mantine/core'

import { AuthentificationSettingsCardWidget } from '@widgets/cerberus-settings/authentification-settings-card/authentification-settings-card.widget'
import { BrandingSettingsCardWidget } from '@widgets/cerberus-settings/branding-settings-card/branding-settings-card.widget'
import { ApiTokensCardWidget } from '@widgets/cerberus-settings/api-tokens-card/api-tokens-card.widget'
import { LoadingScreen, Logo, Page, PageHeaderShared } from '@shared/ui'

interface IProps {
    apiTokensData: FindAllApiTokensCommand.Response['response']
    cerberusSettings: GetCerberusSettingsCommand.Response['response']
}

export const CerberusSettingsPageComponent = (props: IProps) => {
    const { cerberusSettings, apiTokensData } = props

    const { t } = useTranslation()

    if (!cerberusSettings || !apiTokensData) {
        return <LoadingScreen />
    }

    if (
        !cerberusSettings.oauth2Settings ||
        !cerberusSettings.passkeySettings ||
        !cerberusSettings.passwordSettings ||
        !cerberusSettings.tgAuthSettings ||
        !cerberusSettings.brandingSettings
    ) {
        return <LoadingScreen />
    }

    return (
        <Page title={t('constants.cerberus-settings')}>
            <PageHeaderShared icon={<Logo size={24} />} title={t('constants.cerberus-settings')} />
            <Container fluid p={0} size="xl">
                <Masonry columns={{ 300: 1, 1400: 2, 2000: 3, 3000: 4 }} gap={16}>
                    <AuthentificationSettingsCardWidget
                        oauth2Settings={cerberusSettings.oauth2Settings}
                        passkeySettings={cerberusSettings.passkeySettings}
                        passwordSettings={cerberusSettings.passwordSettings}
                        tgAuthSettings={cerberusSettings.tgAuthSettings}
                    />

                    <ApiTokensCardWidget apiTokensData={apiTokensData} />
                    <BrandingSettingsCardWidget
                        brandingSettings={cerberusSettings.brandingSettings}
                    />
                </Masonry>
            </Container>
        </Page>
    )
}
