import NiceModal, { useModal } from '@ebay/nice-modal-react'
import { Modal } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { TbDevices } from 'react-icons/tb'

import { useNiceMantineModal } from '@shared/_modals/use-nice-modal'
import { useIsMobile } from '@shared/hooks/use-is-mobile'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { UserHwidDevicesContentModal } from './user-hwid-devices.content.modal'

interface IProps {
    userId: number
}

export const UserHwidDevicesModal = NiceModal.create((props: IProps) => {
    const { userId } = props
    const { t } = useTranslation()
    const isMobile = useIsMobile()

    const modal = useModal()
    const { modalProps } = useNiceMantineModal({
        modal
    })

    return (
        <Modal
            {...modalProps}
            size="min(1500px, 90vw)"
            styles={{
                body: {
                    height: 'calc(100% - 60px)',
                    display: 'flex',
                    flexDirection: 'column'
                }
            }}
            title={
                <BaseOverlayHeader
                    iconColor="violet"
                    IconComponent={TbDevices}
                    iconVariant="soft"
                    title={t('get-hwid-user-devices.feature.hwid-devices')}
                />
            }
        >
            <UserHwidDevicesContentModal userId={userId} mobile={isMobile} />
        </Modal>
    )
})
