import { Chart } from '@highcharts/react'
import { Box, Card, Center, Group, Loader, SimpleGrid, Stack, Text, ThemeIcon } from '@mantine/core'
import ColorHash from 'color-hash'
import dayjs from 'dayjs'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { PiAppWindowDuotone, PiChartBar } from 'react-icons/pi'

import { useGetSubscriptionRequestHistoryStats } from '@shared/api/hooks'
import { TopLeaderboardCardShared } from '@shared/ui/leaderboard-item-card'
import { formatInt } from '@shared/utils/misc'

export function SrhInspectorMetrics() {
    const { t } = useTranslation()

    const { data: stats, isLoading } = useGetSubscriptionRequestHistoryStats()

    const ch = new ColorHash({ lightness: [0.65, 0.65, 0.65] })

    const getBarChartConfig = () => {
        return {
            chart: {
                type: 'column',
                backgroundColor: 'transparent',
                height: 300,
                style: {
                    fontFamily: 'inherit'
                }
            },
            time: {
                useUTC: false
            },
            accessibility: {
                enabled: false
            },
            title: {
                text: undefined
            },
            xAxis: {
                type: 'datetime',
                dateTimeLabelFormats: {
                    hour: '%H:00'
                },

                units: [['hour', [1, 2, 3, 4, 6, 8, 12]]],
                labels: {
                    style: {
                        color: 'var(--mantine-color-text)'
                    }
                },
                gridLineColor: 'var(--mantine-color-gray-6)',
                lineColor: 'var(--mantine-color-gray-6)'
            },
            yAxis: {
                title: {
                    text: t('srh-inspector-metrics.widget.requests'),
                    style: {
                        color: 'var(--mantine-color-text)'
                    }
                },
                labels: {
                    style: {
                        color: 'var(--mantine-color-text)'
                    }
                },
                gridLineColor: 'var(--mantine-color-gray-6)',
                lineColor: 'var(--mantine-color-gray-6)'
            },

            tooltip: {
                shared: true,
                backgroundColor: 'var(--mantine-color-body)',
                borderColor: 'var(--mantine-color-gray-4)',
                headerFormat: '',
                style: {
                    color: 'var(--mantine-color-text)'
                },
                pointFormatter(this: { x: number; y: number }): string {
                    return `<b>${dayjs(this.x).format('DD.MM.YYYY, HH:mm')}</b> </br> Requests: <b>${this.y}<b>`
                }
            },
            plotOptions: {
                column: {
                    borderWidth: 0,
                    borderRadius: 4,
                    pointPadding: 0.1,
                    groupPadding: 0.1,
                    color: 'var(--mantine-color-indigo-6)',
                    states: {
                        hover: {
                            color: 'var(--mantine-color-indigo-3)'
                        }
                    },
                    dataLabels: {
                        enabled: false
                    }
                }
            },
            legend: {
                enabled: false
            },
            credits: {
                enabled: false
            }
        }
    }

    const hourlyChartOptions = useMemo(() => {
        if (!stats?.hourlyRequestStats || stats.hourlyRequestStats.length === 0) return {}

        const data = stats.hourlyRequestStats.map((item) => {
            const utcDate = new Date(item.dateTime)
            return [utcDate.getTime(), item.requestCount]
        })

        return {
            ...getBarChartConfig(),
            series: [
                {
                    type: 'column',
                    name: t('srh-inspector-metrics.widget.requests'),
                    data
                }
            ]
        }
    }, [stats?.hourlyRequestStats, t])

    const loaderCard = (
        <Center h={300}>
            <Loader size="lg" />
        </Center>
    )

    return (
        <Stack gap="xl">
            <SimpleGrid cols={{ base: 1, lg: 2 }} spacing="lg">
                <Card
                    p="lg"
                    style={{
                        background:
                            'linear-gradient(135deg, var(--mantine-color-dark-6) 0%, var(--mantine-color-dark-7) 100%)'
                    }}
                >
                    <Group align="center" gap="sm" mb="lg" wrap="nowrap">
                        <ThemeIcon color="teal" size="lg" variant="outline">
                            <PiAppWindowDuotone size="20px" />
                        </ThemeIcon>

                        <Text
                            fw={600}
                            size="lg"
                            style={{
                                textOverflow: 'ellipsis',
                                whiteSpace: 'nowrap',
                                overflow: 'hidden'
                            }}
                        >
                            {t('hwid-inspector-metrics.widget.app-distribution')}
                        </Text>
                    </Group>
                    <TopLeaderboardCardShared
                        emptyText={t('hwid-inspector-metrics.widget.no-app-data')}
                        formatValue={(value) => formatInt(value, { thousandSeparator: ' ' })}
                        isLoading={isLoading}
                        items={stats?.byParsedApp?.map((item) => ({
                            color: ch.hex(item.app || 'Unknown'),
                            name: item.app || 'Unknown',
                            total: item.count
                        }))}
                        maxHeight={300}
                        wrapper={(children) => children}
                    />
                </Card>

                <Card
                    p="lg"
                    style={{
                        background:
                            'linear-gradient(135deg, var(--mantine-color-dark-6) 0%, var(--mantine-color-dark-7) 100%)'
                    }}
                >
                    <Group align="center" gap="sm" mb="lg" wrap="nowrap">
                        <ThemeIcon color="blue" size="lg" variant="outline">
                            <PiChartBar size="20px" />
                        </ThemeIcon>

                        <Text
                            fw={600}
                            size="lg"
                            style={{
                                textOverflow: 'ellipsis',
                                whiteSpace: 'nowrap',
                                overflow: 'hidden'
                            }}
                        >
                            {t('srh-inspector-metrics.widget.hourly-request-statistics')}
                        </Text>
                    </Group>
                    {isLoading ? (
                        loaderCard
                    ) : stats?.hourlyRequestStats && stats.hourlyRequestStats.length > 0 ? (
                        <Box h={300}>
                            <Chart options={hourlyChartOptions} />
                        </Box>
                    ) : (
                        <Center h={300}>
                            <Stack align="center" gap="sm">
                                <PiChartBar color="var(--mantine-color-gray-4)" size="48px" />
                                <Text c="dimmed" size="sm">
                                    {t('srh-inspector-metrics.widget.no-hourly-data-available')}
                                </Text>
                            </Stack>
                        </Center>
                    )}
                </Card>
            </SimpleGrid>
        </Stack>
    )
}
