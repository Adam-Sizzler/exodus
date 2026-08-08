import { ActionIcon, Group } from '@mantine/core'
import { GetHostsCommand, GetSubscriptionTemplateCommand } from '@exodus/backend-contract'
import { SubscriptionTemplateEditorWidget } from '@widgets/dashboard/templates/subscription-template-editor'
import { TbArrowBackUp } from 'react-icons/tb'
import { useNavigate } from 'react-router'

import { HelpActionIconShared } from '@shared/_modals/universal'
import { ROUTES } from '@shared/constants'
import { Page, PageHeaderShared } from '@shared/ui'
import { getCoreLogoFromType } from '@shared/ui/get-core-logo-from-type'

interface Props {
    editorType: 'json' | 'yaml'
    hosts: GetHostsCommand.Response['response']
    template: GetSubscriptionTemplateCommand.Response['response']
    title: string
}

export const TemplateEditorPageComponent = (props: Props) => {
    const { editorType, hosts, template, title } = props
    const navigate = useNavigate()

    let isHelpDrawerVisible = false

    if (template.templateType === 'XRAY_JSON') {
        isHelpDrawerVisible = true
    }

    return (
        <Page title={title}>
            <PageHeaderShared
                actions={
                    <Group>
                        <HelpActionIconShared
                            hidden={!isHelpDrawerVisible}
                            screen="EDITOR_TEMPLATES_XRAY_JSON"
                        />

                        <ActionIcon
                            color="gray"
                            onClick={() =>
                                navigate(
                                    ROUTES.DASHBOARD.TEMPLATES.TEMPLATES_BY_TYPE.replace(
                                        ':type',
                                        template.templateType
                                    )
                                )
                            }
                            size="input-md"
                            variant="soft"
                        >
                            <TbArrowBackUp size={24} />
                        </ActionIcon>
                    </Group>
                }
                icon={getCoreLogoFromType({ type: template.templateType })}
                title={template.name}
            />
            <SubscriptionTemplateEditorWidget
                editorType={editorType}
                hosts={hosts}
                template={template}
            />
        </Page>
    )
}
