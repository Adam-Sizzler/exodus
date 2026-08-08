import { useAutoAnimate } from '@formkit/auto-animate/react'
import { Center, RollingNumber, SimpleGrid, Stack, Text } from '@mantine/core'
import { GetHttpStatsCommand } from '@exodus/backend-contract'
import { useTranslation } from 'react-i18next'
import { PiEmpty } from 'react-icons/pi'
import { TbDivide, TbFlame, TbRoute, TbTrendingUp } from 'react-icons/tb'

import { MetricCardShared } from '@shared/ui/metrics/metric-card'
import { SectionCard } from '@shared/ui/section-card'

import { HttpStatRow } from './http-stat-row'

interface IProps {
    httpStats: GetHttpStatsCommand.Response['response']
}

export const HttpStatsTableWidget = ({ httpStats }: IProps) => {
    const { t } = useTranslation()

    const [listRef] = useAutoAnimate<HTMLDivElement>({ duration: 350, easing: 'ease-in-out' })

    const total = httpStats.total
    const avgPerRoute =
        httpStats.routes.length > 0 ? Math.round(total / httpStats.routes.length) : 0

    return (
        <Stack gap="md">
            <SimpleGrid cols={{ base: 1, sm: 2, lg: 4 }} spacing="md">
                <MetricCardShared
                    iconColor="cyan"
                    IconComponent={TbTrendingUp}
                    iconVariant="soft"
                    title={t('http-stats-table.widget.total-requests')}
                    value=""
                    rollingNumberComponent={
                        <RollingNumber value={total} thousandSeparator fw={700} />
                    }
                />
                <MetricCardShared
                    iconColor="blue"
                    IconComponent={TbRoute}
                    iconVariant="soft"
                    title={t('http-stats-table.widget.tracked-routes')}
                    value=""
                    rollingNumberComponent={
                        <RollingNumber value={httpStats.routes.length} thousandSeparator fw={700} />
                    }
                />
                <MetricCardShared
                    iconColor="teal"
                    IconComponent={TbFlame}
                    iconVariant="soft"
                    title={t('http-stats-table.widget.peak-requests')}
                    value=""
                    rollingNumberComponent={
                        <RollingNumber
                            value={httpStats.routes[0]?.count ?? 0}
                            thousandSeparator
                            fw={700}
                        />
                    }
                />
                <MetricCardShared
                    iconColor="grape"
                    IconComponent={TbDivide}
                    iconVariant="soft"
                    title={t('http-stats-table.widget.avg-route')}
                    value=""
                    rollingNumberComponent={
                        <RollingNumber value={avgPerRoute} thousandSeparator fw={700} />
                    }
                />
            </SimpleGrid>

            <SectionCard.Root>
                <SectionCard.Section>
                    {httpStats.routes.length === 0 ? (
                        <Center py="xl">
                            <Stack align="center" gap="xs">
                                <PiEmpty color="var(--mantine-color-gray-5)" size="3rem" />
                                <Text c="dimmed" size="sm">
                                    {t('common.nothing-found')}
                                </Text>
                            </Stack>
                        </Center>
                    ) : (
                        <Stack gap="xs" ref={listRef}>
                            {httpStats.routes.map((route) => (
                                <HttpStatRow
                                    key={`${route.method}|${route.route}`}
                                    percentOfTotal={total > 0 ? (route.count / total) * 100 : 0}
                                    route={route}
                                />
                            ))}
                        </Stack>
                    )}
                </SectionCard.Section>
            </SectionCard.Root>
        </Stack>
    )
}
