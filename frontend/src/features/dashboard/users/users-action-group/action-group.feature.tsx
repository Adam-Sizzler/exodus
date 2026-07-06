import {
    ActionIcon,
    ActionIconGroup,
    Drawer,
    Group,
    Stack,
    Table,
    Text,
    Tooltip
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { modals } from '@mantine/modals'
import { BulkAllUsersActionsWidget } from '@widgets/dashboard/users/bulk-all-users-actions/bulk-all-users-actions.widget'
import { useTranslation } from 'react-i18next'
import {
    TbBaselineDensityLarge,
    TbBaselineDensityMedium,
    TbBaselineDensitySmall,
    TbColumns,
    TbDots,
    TbFilter,
    TbFilterOff,
    TbMaximize,
    TbMinimize,
    TbPlus,
    TbQuestionMark,
    TbRefresh,
    TbRestore,
    TbSettings
} from 'react-icons/tb'

import { QueryKeys } from '@shared/api/hooks'
import { queryClient } from '@shared/api/query-client'
import { useIsMobile } from '@shared/hooks'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { useUserCreationModalStoreActions } from '@entities/dashboard/user-creation-modal-store'
import { useUsersTableStoreActions } from '@entities/dashboard/users/users-table-store'

import { UsersTableTemplatesFeature } from '../users-table-templates/users-table-templates.feature'
import { IProps } from './interfaces'

export const UserActionGroupFeature = (props: IProps) => {
    const { t } = useTranslation()

    const isMobile = useIsMobile()

    const [isHelpDrawerOpen, helpDrawerHandlers] = useDisclosure(false)

    const { isLoading, refetch, table } = props
    const actions = useUsersTableStoreActions()

    const userCreationModalActions = useUserCreationModalStoreActions()

    const handleOpenCreateUserModal = () => {
        userCreationModalActions.changeModalState(true)
    }

    const handleRefetch = () => {
        if (table && refetch) {
            refetch()
        }
    }

    const handleResetTable = () => {
        if (table && refetch) {
            refetch()
            actions.resetState()

            table.resetPageIndex(false)
            table.resetSorting(true)
            table.resetPagination(false)
            table.resetColumnFilters(true)
            table.resetGlobalFilter(true)
        }
    }

    const handleClearFilters = () => {
        if (table && refetch) {
            refetch()

            table.resetPageIndex(false)
            table.resetSorting(true)
            table.resetPagination(false)
            table.resetColumnFilters(true)
            table.resetGlobalFilter(true)
        }
    }

    const handleCloseModal = async () => {
        await queryClient.refetchQueries({ queryKey: QueryKeys.users.getAllUsers._def })
        await queryClient.refetchQueries({ queryKey: QueryKeys.system._def })
    }

    if (!table || !refetch) {
        return null
    }

    return (
        <>
            <Group grow preventGrowOverflow={false} wrap="wrap">
                <ActionIconGroup>
                    <Tooltip label={t('action-group.feature.help')} withArrow>
                        <ActionIcon
                            color="lime"
                            onClick={helpDrawerHandlers.open}
                            size="input-md"
                            variant="soft"
                        >
                            <TbQuestionMark size="24px" />
                        </ActionIcon>
                    </Tooltip>
                </ActionIconGroup>

                <ActionIconGroup>
                    <UsersTableTemplatesFeature table={table} />

                    <Tooltip label={t('action-group.feature.clear-filters')} withArrow>
                        <ActionIcon
                            color="gray"
                            loading={isLoading}
                            onClick={handleClearFilters}
                            size="input-md"
                            variant="soft"
                        >
                            <TbFilterOff size="24px" />
                        </ActionIcon>
                    </Tooltip>

                    <Tooltip label={t('action-group.feature.reset-table')} withArrow>
                        <ActionIcon
                            color="gray"
                            loading={isLoading}
                            onClick={handleResetTable}
                            size="input-md"
                            variant="soft"
                        >
                            <TbRestore size="24px" />
                        </ActionIcon>
                    </Tooltip>
                </ActionIconGroup>

                <ActionIconGroup>
                    <Tooltip label={t('action-group.feature.bulk-actions')} withArrow>
                        <ActionIcon
                            color="red"
                            loading={isLoading}
                            onClick={() =>
                                modals.open({
                                    title: (
                                        <BaseOverlayHeader
                                            iconColor="cyan"
                                            IconComponent={TbDots}
                                            iconVariant="soft"
                                            title={t(
                                                'bulk-all-user-actions-drawer.widget.bulk-all-user-actions'
                                            )}
                                            titleOrder={5}
                                        />
                                    ),
                                    onClose: handleCloseModal,
                                    centered: true,
                                    size: 'lg',
                                    fullScreen: isMobile,
                                    children: <BulkAllUsersActionsWidget isMobile={isMobile} />
                                })
                            }
                            size="input-md"
                            variant="soft"
                        >
                            <TbSettings size="24px" />
                        </ActionIcon>
                    </Tooltip>

                    <Tooltip label={t('common.update')} withArrow>
                        <ActionIcon
                            loading={isLoading}
                            onClick={handleRefetch}
                            size="input-md"
                            variant="soft"
                        >
                            <TbRefresh size="24px" />
                        </ActionIcon>
                    </Tooltip>

                    <Tooltip label={t('action-group.feature.new-user')} withArrow>
                        <ActionIcon
                            color="teal"
                            onClick={handleOpenCreateUserModal}
                            size="input-md"
                            variant="soft"
                        >
                            <TbPlus size="24px" />
                        </ActionIcon>
                    </Tooltip>
                </ActionIconGroup>
            </Group>

            <Drawer
                keepMounted={false}
                onClose={helpDrawerHandlers.close}
                opened={isHelpDrawerOpen}
                overlayProps={{ backgroundOpacity: 0.6, blur: 0 }}
                padding="lg"
                position="left"
                size="500px"
                title={
                    <BaseOverlayHeader
                        iconColor="teal"
                        IconComponent={TbQuestionMark}
                        iconVariant="soft"
                        title={t('action-group.feature.table-controls-help')}
                    />
                }
            >
                <Stack gap="md">
                    <Text c="dimmed" size="sm">
                        {t('action-group.feature.table-controler-description-line-1')}
                    </Text>

                    <Table>
                        <Table.Thead>
                            <Table.Tr>
                                <Table.Th>
                                    {t('action-group.feature.table-controler-description-line-2')}
                                </Table.Th>
                                <Table.Th>
                                    {t('action-group.feature.table-controler-description-line-3')}
                                </Table.Th>
                            </Table.Tr>
                        </Table.Thead>
                        <Table.Tbody>
                            <Table.Tr>
                                <Table.Td>
                                    <Group gap="xs">
                                        <TbFilter size="24px" />
                                        <TbFilterOff size="24px" />
                                    </Group>
                                </Table.Td>
                                <Table.Td>
                                    {t('action-group.feature.table-controler-description-line-4')}
                                </Table.Td>
                            </Table.Tr>
                            <Table.Tr>
                                <Table.Td>
                                    <TbColumns size="24px" />
                                </Table.Td>
                                <Table.Td>
                                    {t('action-group.feature.show-or-hide-specific-columns')}
                                </Table.Td>
                            </Table.Tr>
                            <Table.Tr>
                                <Table.Td>
                                    <Group gap="xs">
                                        <TbBaselineDensitySmall size="24px" />
                                        <TbBaselineDensityMedium size="24px" />
                                        <TbBaselineDensityLarge size="24px" />
                                    </Group>
                                </Table.Td>
                                <Table.Td>
                                    {t('action-group.feature.adjust-row-spacing-density')}
                                </Table.Td>
                            </Table.Tr>
                            <Table.Tr>
                                <Table.Td>
                                    <Group gap="xs">
                                        <TbMaximize size="24px" />
                                        <TbMinimize size="24px" />
                                    </Group>
                                </Table.Td>
                                <Table.Td>
                                    {t('action-group.feature.toggle-fullscreen-table-view')}
                                </Table.Td>
                            </Table.Tr>
                        </Table.Tbody>
                    </Table>
                </Stack>
            </Drawer>
        </>
    )
}
