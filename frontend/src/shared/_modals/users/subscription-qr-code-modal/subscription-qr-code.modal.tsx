import NiceModal, { useModal } from '@ebay/nice-modal-react'
import { Modal } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { TbQrcode } from 'react-icons/tb'

import { useNiceMantineModal } from '@shared/_modals/use-nice-modal'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { QrCodeBuilder } from '@shared/ui/qr-code-builder'

interface IProps {
    subscriptionUrl: string
    username: string
}

export const SubscriptionQrCodeModal = NiceModal.create((props: IProps) => {
    const { subscriptionUrl, username } = props
    const { t } = useTranslation()

    const modal = useModal()
    const { modalProps } = useNiceMantineModal({
        modal
    })

    return (
        <Modal
            {...modalProps}
            size="auto"
            title={
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={TbQrcode}
                    iconVariant="soft"
                    title={t('view-user-modal.widget.subscription-qr-code')}
                />
            }
        >
            <QrCodeBuilder data={subscriptionUrl} title={username} />
        </Modal>
    )
})
