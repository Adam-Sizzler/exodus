import {
    ActionIcon,
    Box,
    Card,
    Group,
    Loader,
    Stack,
    Text,
    ThemeIcon,
    Tooltip
} from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { TbDevices, TbRefresh } from 'react-icons/tb'

import { useGetUserSubscriptionRequestHistory } from '@shared/api/hooks'
import { EmptyPageLayout } from '@shared/ui/layouts/empty-page'
import { LoaderModalShared } from '@shared/ui/loader-modal'

import { UserSubscriptionRequestItem } from './user-subscription-request-item'
import classes from './user-subscription-requests.module.css'
import { UserSubscriptionRequestsTable } from './user-subscription-requests.table'

interface IProps {
    mobile: boolean
    userUuid: string
}

export const UserSubscriptionRequestsContentWidget = (props: IProps) => {
    const { userUuid, mobile } = props
    const { t } = useTranslation()

    const {
        data: records,
        isLoading,
        isRefetching,
        refetch
    } = useGetUserSubscriptionRequestHistory({
        route: {
            uuid: userUuid
        }
    })

    return (
        <Stack className={classes.drawerContent}>
            <Card withBorder>
                <Stack gap="md">
                    <Group gap="sm" justify="space-between">
                        <Group>
                            <ThemeIcon color="indigo" radius="md" size="xl" variant="soft">
                                <TbDevices size={24} />
                            </ThemeIcon>
                            <Stack gap={0}>
                                <Box
                                    style={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        height: 'calc(var(--mantine-font-size-xl) * var(--mantine-line-height))'
                                    }}
                                >
                                    {isLoading ? (
                                        <Loader color="cyan" size="sm" type="oval" />
                                    ) : (
                                        <Text c="white" fw={700} size="xl">
                                            {records?.total ?? 0}
                                        </Text>
                                    )}
                                </Box>
                                <Text c="dimmed" size="xs">
                                    {t('user-subscription-requests-drawer.widget.total-records')}
                                </Text>
                            </Stack>
                        </Group>
                        <Group gap="xs">
                            <Tooltip label={t('common.refresh')}>
                                <ActionIcon
                                    color="indigo"
                                    loading={isRefetching}
                                    onClick={() => refetch()}
                                    size="lg"
                                    variant="soft"
                                >
                                    <TbRefresh size={20} />
                                </ActionIcon>
                            </Tooltip>
                        </Group>
                    </Group>
                </Stack>
            </Card>

            {mobile && !isLoading && records?.total === 0 && <EmptyPageLayout />}

            {mobile && isLoading && <LoaderModalShared h="80vh" text="Loading..." w="100%" />}
            {mobile && !isLoading && records && records.records.length > 0 && (
                <Stack>
                    {records.records.map((record) => (
                        <UserSubscriptionRequestItem key={record.id} request={record} />
                    ))}
                </Stack>
            )}

            {!mobile && (
                <Box className={classes.tableContainer}>
                    <UserSubscriptionRequestsTable
                        isLoading={isLoading}
                        records={records?.records}
                    />
                </Box>
            )}
        </Stack>
    )
}
