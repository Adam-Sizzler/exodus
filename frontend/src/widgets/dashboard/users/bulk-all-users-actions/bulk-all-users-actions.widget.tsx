import { DeleteAllUsersByStatusFeature } from '@features/ui/dashboard/users/delete-all-users-by-status'
import { Button, NumberInput, Stack } from '@mantine/core'
import { modals } from '@mantine/modals'
import { useTranslation } from 'react-i18next'
import { TbCalendarUp, TbPencil, TbRefresh, TbTrash, TbUsers, TbUsersMinus } from 'react-icons/tb'

import { useBulkAllExtendUsersExpirationDate, useBulkAllResetTrafficUsers } from '@shared/api/hooks'
import { ActionCardShared } from '@shared/ui'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { BulkAllUsersUpdateWidget } from './bulk-all-users-update.widget'

interface IProps {
    isMobile: boolean
}

export const BulkAllUsersActionsWidget = (props: IProps) => {
    const { isMobile } = props

    const { t } = useTranslation()

    const { mutate: resetTraffic, isPending: isResetTrafficPending } = useBulkAllResetTrafficUsers({
        mutationFns: {}
    })

    const { mutate: extendExpirationDate, isPending: isExtendExpirationDatePending } =
        useBulkAllExtendUsersExpirationDate({
            mutationFns: {}
        })

    const handleExtendExpirationDate = () => {
        let userInput = 1

        modals.open({
            title: (
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={TbCalendarUp}
                    iconVariant="soft"
                    title={t('bulk-user-actions.actions.tab.feature.extend-expiration-date')}
                    titleOrder={5}
                />
            ),
            centered: true,
            children: (
                <>
                    <NumberInput
                        allowDecimal={false}
                        allowNegative={false}
                        data-autofocus
                        decimalScale={0}
                        defaultValue={1}
                        description={t(
                            'bulk-user-actions.actions.tab.feature.enter-the-number-of-days-to-extend-the-expiration-date'
                        )}
                        label={t('bulk-user-actions.actions.tab.feature.extend-days')}
                        max={9999}
                        min={1}
                        onChange={(value) => {
                            userInput = Number(value)
                        }}
                        required
                        step={1}
                    />
                    <Button
                        fullWidth
                        mt="md"
                        onClick={() => {
                            modals.closeAll()
                            extendExpirationDate({
                                variables: {
                                    extendDays: userInput
                                }
                            })
                        }}
                    >
                        {t('bulk-user-actions.actions.tab.feature.extend')}
                    </Button>
                </>
            )
        })
    }

    return (
        <Stack gap="xs">
            <ActionCardShared
                description={t('bulk-all-user-actions-drawer.widget.all-users')}
                icon={<TbPencil size={20} />}
                iconColor="teal"
                onClick={() => {
                    modals.open({
                        title: (
                            <BaseOverlayHeader
                                iconColor="cyan"
                                IconComponent={TbUsers}
                                iconVariant="soft"
                                title={t(
                                    'bulk-all-user-actions-tabs.update.tab.feature.update-users'
                                )}
                            />
                        ),
                        size: '1000px',
                        fullScreen: isMobile,
                        centered: true,
                        children: <BulkAllUsersUpdateWidget isMobile={isMobile} />
                    })
                }}
                title={t('common.update')}
                variant="soft"
            />
            <ActionCardShared
                description={t(
                    'bulk-user-actions.actions.tab.feature.extend-expiration-date-description'
                )}
                icon={<TbCalendarUp size={20} />}
                iconColor="cyan"
                isLoading={isExtendExpirationDatePending}
                onClick={handleExtendExpirationDate}
                title={t('bulk-user-actions.actions.tab.feature.extend-expiration-date')}
                variant="soft"
            />
            <ActionCardShared
                description={t(
                    'bulk-all-user-actions-tabs.actions.tab.feature.reset-traffic-description'
                )}
                icon={<TbRefresh size={20} />}
                iconColor="blue"
                isLoading={isResetTrafficPending}
                onClick={() => resetTraffic({})}
                title={t('bulk-all-user-actions-tabs.actions.tab.feature.reset-traffic')}
                variant="soft"
                withConfirmation
            />
            <ActionCardShared
                description={t(
                    'bulk-all-user-actions-tabs.danger.tab.feature.delete-users-by-status-description'
                )}
                icon={<TbUsersMinus size={20} />}
                iconColor="red"
                onClick={() => {
                    modals.open({
                        title: (
                            <BaseOverlayHeader
                                iconColor="red"
                                IconComponent={TbTrash}
                                iconVariant="soft"
                                title={t(
                                    'bulk-all-user-actions-tabs.danger.tab.feature.delete-users-by-status'
                                )}
                            />
                        ),
                        centered: true,
                        children: <DeleteAllUsersByStatusFeature />
                    })
                }}
                title={t('bulk-all-user-actions-tabs.danger.tab.feature.delete-users-by-status')}
                variant="soft"
            />
        </Stack>
    )
}
