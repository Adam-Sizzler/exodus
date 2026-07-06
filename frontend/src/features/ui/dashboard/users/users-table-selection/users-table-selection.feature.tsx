import { Button, Group, Text } from '@mantine/core'
import { modals } from '@mantine/modals'
import { BulkUsersActionsWidget } from '@widgets/dashboard/users/bulk-users-actions/bulk-users-actions.widget'
import { useTranslation } from 'react-i18next'
import { PiClockClockwise } from 'react-icons/pi'
import { TbDots } from 'react-icons/tb'

import { QueryKeys } from '@shared/api/hooks'
import { queryClient } from '@shared/api/query-client'
import { useIsMobile } from '@shared/hooks'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { useBulkUsersActionsStoreActions } from '@entities/dashboard/users/bulk-users-actions-store'

import { IProps } from './interfaces/props.interface'

export const UsersTableSelectionFeature = (props: IProps) => {
    const { resetRowSelection, toggleAllPageRowsSelected } = props
    const { t } = useTranslation()
    const isMobile = useIsMobile()

    const bulkUsersActionsStoreActions = useBulkUsersActionsStoreActions()

    const handleClearSelection = () => {
        resetRowSelection()
        bulkUsersActionsStoreActions.resetState()
    }

    const usersToUpdate = bulkUsersActionsStoreActions.getUuidLength()

    if (usersToUpdate === 0) {
        return null
    }

    const handleCloseModal = async () => {
        resetRowSelection()
        bulkUsersActionsStoreActions.resetState()
        await queryClient.refetchQueries({ queryKey: QueryKeys.users.getAllUsers._def })
        await queryClient.refetchQueries({ queryKey: QueryKeys.system._def })
    }

    return (
        <Group justify="apart" px="xs">
            <Text fw={600} size="sm">
                {usersToUpdate} {t('users-table-selection.feature.row-s-selected')}
            </Text>
            <Group gap="xs">
                <Button color="blue" onClick={handleClearSelection} size="xs" variant="subtle">
                    {t('users-table-selection.feature.clear-selection')}
                </Button>
                <Button color="blue" onClick={toggleAllPageRowsSelected} size="xs" variant="subtle">
                    {t('users-table-selection.feature.select-all')}
                </Button>
                <Button
                    color="green"
                    leftSection={<PiClockClockwise />}
                    onClick={() =>
                        modals.open({
                            title: (
                                <BaseOverlayHeader
                                    iconColor="cyan"
                                    IconComponent={TbDots}
                                    iconVariant="soft"
                                    subtitle={t(
                                        'bulk-user-actions.actions.tab.feature.perform-action-on-users',
                                        {
                                            usersCount: usersToUpdate
                                        }
                                    )}
                                    title={t('bulk-user-actions-drawer.widget.bulk-user-actions')}
                                    titleOrder={5}
                                />
                            ),
                            size: 'lg',
                            fullScreen: isMobile,
                            centered: true,
                            onClose: handleCloseModal,
                            children: <BulkUsersActionsWidget isMobile={isMobile} />
                        })
                    }
                    size="sm"
                    variant="subtle"
                >
                    {t('users-table-selection.feature.bulk-actions')}
                </Button>
            </Group>
        </Group>
    )
}
