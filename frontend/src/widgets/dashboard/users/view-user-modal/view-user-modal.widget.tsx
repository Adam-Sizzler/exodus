import { Modal } from '@mantine/core'
import { IconUser } from '@tabler/icons-react'
import { motion } from 'motion/react'
import { useTranslation } from 'react-i18next'

import { queryClient } from '@shared/api'
import { usersQueryKeys } from '@shared/api/hooks'
import { useIsMobile } from '@shared/hooks'
import { LoaderModalShared } from '@shared/ui/loader-modal'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import {
    useUserModalStoreActions,
    useUserModalStoreIsModalOpen,
    useUserModalStoreUserUuid
} from '@entities/dashboard/user-modal-store/user-modal-store'
import { useBulkUsersActionsStoreActions } from '@entities/dashboard/users/bulk-users-actions-store'

import { ViewUserModalContent } from './view-user-modal.content'

export const ViewUserModal = () => {
    const { t } = useTranslation()

    const isViewUserModalOpen = useUserModalStoreIsModalOpen()
    const actions = useUserModalStoreActions()
    const bulkUsersActionsStoreActions = useBulkUsersActionsStoreActions()
    const selectedUser = useUserModalStoreUserUuid()

    const isMobile = useIsMobile()

    const handleClose = () => {
        bulkUsersActionsStoreActions.resetState()
        actions.clearModalState()

        if (selectedUser) {
            queryClient.removeQueries({
                queryKey: usersQueryKeys.getUserByUuid({
                    uuid: selectedUser
                }).queryKey
            })
        }
    }

    return (
        <Modal
            centered
            fullScreen={isMobile}
            onClose={() => actions.changeModalState(false)}
            onExitTransitionEnd={handleClose}
            opened={isViewUserModalOpen}
            size="1000px"
            title={
                <BaseOverlayHeader
                    iconColor="cyan"
                    IconComponent={IconUser}
                    iconVariant="soft"
                    title={t('view-user-modal.widget.edit-user-headline')}
                />
            }
            transitionProps={isMobile ? { transition: 'fade', duration: 200 } : undefined}
        >
            {!selectedUser ? (
                <motion.div
                    animate={{ opacity: 1 }}
                    initial={{ opacity: 0 }}
                    transition={{ duration: 0.3 }}
                >
                    <LoaderModalShared h="78vh" />
                </motion.div>
            ) : (
                <ViewUserModalContent userUuid={selectedUser} />
            )}
        </Modal>
    )
}
