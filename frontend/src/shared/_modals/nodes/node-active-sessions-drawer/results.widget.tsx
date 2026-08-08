import { Center, Stack, Tabs, Text } from '@mantine/core'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { PiEmptyDuotone } from 'react-icons/pi'
import { TbSortAscending, TbSortDescending } from 'react-icons/tb'

import { SectionCard } from '@shared/ui/section-card'

import { IpStatsWidget } from './ip-stats.widget'
import { SessionsListWidget } from './sessions-list.widget'
import { SummaryCardWidget } from './summary-card.widget'
import { ActiveSessionUser } from './use-node-active-sessions'

type SortMode = 'default' | 'ips-asc' | 'ips-desc'

interface IProps {
    onRefresh: () => void
    users: ActiveSessionUser[]
}

export const ResultsWidget = ({ onRefresh, users }: IProps) => {
    const { t } = useTranslation()
    const [sortMode, setSortMode] = useState<SortMode>('default')

    const sortedUsers = useMemo(() => {
        switch (sortMode) {
            case 'ips-asc':
                return [...users].sort((a, b) => a.ips.length - b.ips.length)
            case 'ips-desc':
                return [...users].sort((a, b) => b.ips.length - a.ips.length)
            default:
                return users
        }
    }, [users, sortMode])

    const hasUsers = sortedUsers.length > 0

    return (
        <Stack gap="md" style={{ flex: 1, minHeight: 0 }}>
            <SummaryCardWidget
                onRefresh={onRefresh}
                totalUsers={sortedUsers.length}
                ipStats={hasUsers && <IpStatsWidget users={sortedUsers} />}
            />

            {hasUsers && (
                <Tabs
                    onChange={(value) => setSortMode((value as SortMode) ?? 'default')}
                    value={sortMode}
                >
                    <Tabs.List grow>
                        <Tabs.Tab leftSection={<TbSortAscending size={16} />} value="default">
                            Default
                        </Tabs.Tab>
                        <Tabs.Tab leftSection={<TbSortAscending size={16} />} value="ips-asc">
                            IPs
                        </Tabs.Tab>
                        <Tabs.Tab leftSection={<TbSortDescending size={16} />} value="ips-desc">
                            IPs
                        </Tabs.Tab>
                    </Tabs.List>
                </Tabs>
            )}

            {!hasUsers && (
                <SectionCard.Root gap="sm">
                    <SectionCard.Section>
                        <Center h="230">
                            <Stack align="center" gap="xs">
                                <PiEmptyDuotone color="var(--mantine-color-gray-5)" size="3rem" />
                                <Text c="dimmed" size="sm">
                                    {t('active-sessions-drawer.widget.no-active-sessions')}
                                </Text>
                            </Stack>
                        </Center>
                    </SectionCard.Section>
                </SectionCard.Root>
            )}

            {hasUsers && <SessionsListWidget users={sortedUsers} />}
        </Stack>
    )
}
