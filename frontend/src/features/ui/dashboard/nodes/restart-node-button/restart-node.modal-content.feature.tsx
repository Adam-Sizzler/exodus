import { Stack } from '@mantine/core'
import { modals } from '@mantine/modals'
import { useTranslation } from 'react-i18next'
import { TbAlertCircle, TbRocket } from 'react-icons/tb'

import { useRestartNode } from '@shared/api/hooks'
import { ActionCardShared } from '@shared/ui/action-card'

interface IProps {
    handleClose: () => void
    nodeUuid: string
}

export function RestartNodeModalContentFeature(props: IProps) {
    const { handleClose, nodeUuid } = props
    const { t } = useTranslation()

    const mutationParams = {
        route: {
            uuid: nodeUuid
        },
        mutationFns: {
            onSuccess: async () => {
                modals.closeAll()
                handleClose()
            }
        }
    }

    const { mutate: restartNode, isPending } = useRestartNode(mutationParams)

    return (
        <Stack gap="sm">
            <ActionCardShared
                description={t('restart-node-button.feature.restart-description-line-1')}
                icon={<TbAlertCircle size={22} />}
                iconColor="red"
                isLoading={isPending}
                onClick={() => {
                    restartNode({
                        variables: {
                            forceRestart: true
                        }
                    })
                }}
                title={t('nodes-header-action-buttons.feature.force')}
                variant="soft"
            />

            <ActionCardShared
                description={t('restart-node-button.feature.restart-description-line-2')}
                icon={<TbRocket size={22} />}
                iconColor="teal"
                isLoading={isPending}
                onClick={() => {
                    restartNode({
                        variables: {
                            forceRestart: false
                        }
                    })
                }}
                title={t('nodes-header-action-buttons.feature.graceful')}
                variant="soft"
            />
        </Stack>
    )
}
