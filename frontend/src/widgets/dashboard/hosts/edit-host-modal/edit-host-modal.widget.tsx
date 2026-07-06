import { Drawer } from '@mantine/core'
import { memo } from 'react'
import { useTranslation } from 'react-i18next'
import { PiListChecks } from 'react-icons/pi'

import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { MODALS, useModalClose, useModalState } from '@entities/dashboard/modal-store'

import { EditHostModalContentWidget } from './edit-host-modal-content.widget'

export const EditHostModalWidget = memo(() => {
    const { t } = useTranslation()

    const { isOpen, internalState: host } = useModalState(MODALS.EDIT_HOST_MODAL)
    const close = useModalClose(MODALS.EDIT_HOST_MODAL)

    return (
        <Drawer
            keepMounted={false}
            onClose={close}
            opened={isOpen}
            overlayProps={{ backgroundOpacity: 0.6, blur: 0 }}
            padding="lg"
            position="right"
            size="lg"
            title={
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={PiListChecks}
                    iconVariant="soft"
                    subtitle={host?.uuid}
                    title={t('edit-host-modal.widget.edit-host')}
                    withCopy={true}
                />
            }
        >
            {host && <EditHostModalContentWidget host={host} onClose={close} />}
        </Drawer>
    )
})
