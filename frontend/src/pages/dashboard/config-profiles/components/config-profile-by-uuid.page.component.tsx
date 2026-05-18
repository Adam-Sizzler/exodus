import { GetConfigProfileByUuidCommand } from '@exodus/backend-contract'
import { ActionIcon, Group } from '@mantine/core'
import { TbArrowBackUp, TbFile } from 'react-icons/tb'
import { useTranslation } from 'node_modules/react-i18next'
import { useNavigate } from 'react-router-dom'

import { ConfigEditorWidget } from '@widgets/dashboard/config-profiles/config-editor/config-editor.widget'
import { PageHeaderShared } from '@shared/ui/page-header/page-header.shared'
import { HelpActionIconShared } from '@shared/ui/help-drawer'
import { ROUTES } from '@shared/constants'
import { Page } from '@shared/ui/page'

interface Props {
    configProfile: GetConfigProfileByUuidCommand.Response['response']
}

export const ConfigProfileByUuidPageComponent = (props: Props) => {
    const { configProfile } = props

    const { t } = useTranslation()
    const navigate = useNavigate()

    return (
        <>
            <Page title={t('constants.config-profiles')}>
                <PageHeaderShared
                    actions={
                        <Group>
                            <HelpActionIconShared hidden={false} screen="PAGE_CONFIG_PROFILES" />

                            <ActionIcon
                                color="gray"
                                onClick={() =>
                                    navigate(ROUTES.DASHBOARD.MANAGEMENT.CONFIG_PROFILES)
                                }
                                size="input-md"
                                variant="light"
                            >
                                <TbArrowBackUp size={24} />
                            </ActionIcon>
                        </Group>
                    }
                    description={configProfile.uuid}
                    icon={<TbFile size={24} />}
                    title={configProfile.name}
                />

                <ConfigEditorWidget configProfile={configProfile} />
            </Page>
        </>
    )
}
