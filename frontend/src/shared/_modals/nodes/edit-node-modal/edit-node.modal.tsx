import NiceModal, { useModal } from '@ebay/nice-modal-react'
import { Modal } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { TbCpu } from 'react-icons/tb'

import { useNiceMantineModal } from '@shared/_modals/use-nice-modal'
import { queryClient } from '@shared/api'
import { QueryKeys } from '@shared/api/hooks'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { EditNodeByUuidModalContent } from './edit-node.modal.content'

interface IProps {
    nodeUuid: string
}

export const EditNodeModal = NiceModal.create((props: IProps) => {
    const { nodeUuid } = props

    const modal = useModal()
    const { modalProps, hide } = useNiceMantineModal({
        modal,
        onClose() {
            queryClient.refetchQueries({
                queryKey: QueryKeys.nodes.getAllNodes.queryKey
            })
        }
    })

    const { t } = useTranslation()

    return (
        <Modal
            {...modalProps}
            size="1000px"
            title={
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={TbCpu}
                    iconVariant="soft"
                    title={t('edit-node-modal.widget.edit-node')}
                />
            }
        >
            <EditNodeByUuidModalContent nodeUuid={nodeUuid} onClose={hide} />
        </Modal>
    )
})
