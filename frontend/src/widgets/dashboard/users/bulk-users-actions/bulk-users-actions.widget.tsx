import { Button, Divider, Group, NumberInput, Paper, Stack, Text, ThemeIcon } from '@mantine/core'
import { useForm } from '@mantine/form'
import { modals } from '@mantine/modals'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import {
    TbCalendarUp,
    TbCirclesRelation,
    TbPencil,
    TbRefresh,
    TbRefreshAlert,
    TbUsers,
    TbUsersMinus
} from 'react-icons/tb'

import {
    useBulkDeleteUsers,
    useBulkExtendUsersExpirationDate,
    useBulkResetTraffic,
    useBulkRevokeUsersSubscription,
    useBulkSetActiveInternalSquads,
    useGetInternalSquads
} from '@shared/api/hooks'
import { ActionCardShared } from '@shared/ui'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { useBulkUsersActionsStoreActions } from '@entities/dashboard/users/bulk-users-actions-store'

import { InternalSquadsListWidget } from '../internal-squads-list'
import { BulkUsersUpdateWidget } from './bulk-users-update.widget'

interface IProps {
    isMobile: boolean
}

export const BulkUsersActionsWidget = (props: IProps) => {
    const { isMobile } = props
    const { t } = useTranslation()

    const [searchQuery, setSearchQuery] = useState('')

    const actions = useBulkUsersActionsStoreActions()

    const { data: internalSquads } = useGetInternalSquads()

    const { mutate: deleteUsers, isPending: isDeletePending } = useBulkDeleteUsers()
    const { mutate: revokeUsersSubscription, isPending: isRevokePending } =
        useBulkRevokeUsersSubscription()
    const { mutate: resetTraffic, isPending: isResetPending } = useBulkResetTraffic()
    const { mutate: setActiveInternalSquads, isPending: isSetActiveInternalSquadsPending } =
        useBulkSetActiveInternalSquads()
    const { mutate: extendExpirationDate, isPending: isExtendExpirationDatePending } =
        useBulkExtendUsersExpirationDate()

    const uuids = actions.getUuids()

    const form = useForm({
        name: 'change-active-internal-squads-form',
        mode: 'uncontrolled',
        initialValues: {
            activeInternalSquads: []
        }
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
                                    uuids,
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

    const filteredInternalSquads = useMemo(() => {
        const allInternalSquads = internalSquads?.internalSquads || []
        if (!searchQuery.trim()) return allInternalSquads

        const query = searchQuery.toLowerCase().trim()
        return allInternalSquads.filter((internalSquad) =>
            internalSquad.name?.toLowerCase().includes(query)
        )
    }, [internalSquads, searchQuery])

    return (
        <Stack gap="xs">
            <ActionCardShared
                description=""
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
                        children: <BulkUsersUpdateWidget isMobile={isMobile} />
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
                isLoading={isResetPending}
                onClick={() =>
                    resetTraffic({
                        variables: { uuids }
                    })
                }
                title={t('bulk-all-user-actions-tabs.actions.tab.feature.reset-traffic')}
                variant="soft"
                withConfirmation
            />
            <ActionCardShared
                description={t(
                    'bulk-user-actions.actions.tab.feature.revokes-subscription-for-all-selected-users'
                )}
                icon={<TbRefreshAlert size={20} />}
                iconColor="orange"
                isLoading={isRevokePending}
                onClick={() =>
                    revokeUsersSubscription({
                        variables: { uuids }
                    })
                }
                title={t('bulk-user-actions.actions.tab.feature.revoke-subscription')}
                variant="soft"
                withConfirmation
            />
            <ActionCardShared
                description={t(
                    'bulk-user-actions.danger.tab.feature.permanently-deletes-all-selected-users-and-their-data'
                )}
                icon={<TbUsersMinus size={20} />}
                iconColor="red"
                isLoading={isDeletePending}
                onClick={() =>
                    deleteUsers({
                        variables: { uuids }
                    })
                }
                title={t('bulk-user-actions.danger.tab.feature.delete-users')}
                variant="soft"
                withConfirmation
            />
            <Paper
                bg="rgba(255, 255, 255, 0.02)"
                p="md"
                shadow="sm"
                style={{
                    border: '1px solid rgba(255, 255, 255, 0.08)'
                }}
                withBorder
            >
                <Stack gap="md">
                    <Group gap="md" justify="space-between" wrap="nowrap">
                        <Group gap="md" wrap="nowrap">
                            <ThemeIcon color="cyan" radius="md" size="xl" variant="soft">
                                <TbCirclesRelation size={20} />
                            </ThemeIcon>
                            <Stack gap={2}>
                                <Text fw={600} size="sm">
                                    {t('bulk-user-actions.actions.tab.feature.internal-squads')}
                                </Text>
                                <Text c="dimmed" size="xs">
                                    {t(
                                        'bulk-user-actions.actions.tab.feature.specify-internal-squads-that-will-be-assigned-to-the-user'
                                    )}
                                </Text>
                            </Stack>
                        </Group>
                    </Group>

                    <InternalSquadsListWidget
                        filteredInternalSquads={filteredInternalSquads}
                        formKey={form.key('activeInternalSquads')}
                        hideEditButton={true}
                        searchQuery={searchQuery}
                        setSearchQuery={setSearchQuery}
                        {...form.getInputProps('activeInternalSquads')}
                    />
                </Stack>

                <Divider mb="md" mt="xs" />

                <Group justify="flex-end">
                    <Button
                        color="cyan"
                        loading={isSetActiveInternalSquadsPending}
                        onClick={() => {
                            setActiveInternalSquads({
                                variables: {
                                    uuids,
                                    activeInternalSquads: form.getValues().activeInternalSquads
                                }
                            })
                        }}
                        variant="soft"
                    >
                        {t('common.change')}
                    </Button>
                </Group>
            </Paper>
        </Stack>
    )
}
