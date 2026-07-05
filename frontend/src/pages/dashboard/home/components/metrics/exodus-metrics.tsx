import {
    PiClockDuotone,
    PiCpuDuotone,
    PiGearSixDuotone,
    PiMemoryDuotone,
    PiTimerDuotone
} from 'react-icons/pi'
import { TFunction } from 'i18next'

import { IMetricCardProps } from '@shared/ui/metrics/metric-card'
import { prettyBytesToAnyUtil } from '@shared/utils/bytes'

export type ExodusSummary = {
    startedAt?: string
    uptimeSeconds?: number
    totalRssBytes?: number
    heapAllocBytes?: number
    averageCpuPercent?: number
    averageGoroutines?: number
    averageSchedulerDelayMs?: number
    averageSchedulerP99Ms?: number
}

export type ExodusMetric = {
    name?: string
    instanceType?: string
    instanceId?: number
    pid?: number
    startedAt?: string
    uptimeSeconds?: number
    runtime?: {
        language?: string
        version?: string
        goos?: string
        goarch?: string
    }
    cpu?: {
        cores?: number
        gomaxprocs?: number
        processPercent?: number
    }
    memory?: {
        rssBytes?: number
        allocBytes?: number
        sysBytes?: number
        heapAllocBytes?: number
        heapSysBytes?: number
        stackInuseBytes?: number
        heapUsedPercent?: number
    }
    gc?: {
        numGc?: number
    }
    scheduler?: {
        goroutines?: number
        schedulerDelayMs?: number
        schedulerP99Ms?: number
    }
    process?: {
        openFileDescriptors?: number
        threads?: number
    }
}

export const getFirstExodusMetric = (
    exodusMetrics?: ExodusMetric[]
): ExodusMetric | undefined => {
    if (!exodusMetrics || exodusMetrics.length === 0) {
        return undefined
    }

    return exodusMetrics[0]
}

export const formatMilliseconds = (value?: number): string => `${Number(value ?? 0).toFixed(3)} ms`
const formatPercent = (value?: number): string => `${Number(value ?? 0).toFixed(1)} %`
export const formatBytes = (value?: number): string =>
    prettyBytesToAnyUtil(Number(value ?? 0), true)

const getUptimeSeconds = (startedAt?: string, fallbackSeconds?: number): number => {
    if (startedAt) {
        const startedAtMs = Date.parse(startedAt)

        if (Number.isFinite(startedAtMs)) {
            return Math.max(0, Math.floor((Date.now() - startedAtMs) / 1000))
        }
    }

    return Math.max(0, Math.floor(Number(fallbackSeconds ?? 0)))
}

const formatCompactDuration = (uptimeInSeconds: number | undefined, t: TFunction): string => {
    const totalSeconds = Math.max(0, Math.floor(Number(uptimeInSeconds ?? 0)))
    const minute = 60
    const hour = 3_600
    const day = 86_400

    const days = Math.floor(totalSeconds / day)
    const hours = Math.floor((totalSeconds % day) / hour)
    const minutes = Math.floor((totalSeconds % hour) / minute)
    const seconds = totalSeconds % minute

    if (days >= 30) {
        return `${days} ${t('runtime-metrics.duration.days')}`
    }

    if (days > 0) {
        return `${days} ${t('runtime-metrics.duration.days')} ${hours} ${t('runtime-metrics.duration.hours')}`
    }

    if (hours > 0) {
        return `${hours} ${t('runtime-metrics.duration.hours')} ${minutes} ${t('runtime-metrics.duration.minutes')}`
    }

    if (minutes > 0) {
        return `${minutes} ${t('runtime-metrics.duration.minutes')} ${seconds} ${t('runtime-metrics.duration.seconds')}`
    }

    return `${seconds} ${t('runtime-metrics.duration.seconds')}`
}

export const getHeapUsedValue = (metric?: ExodusMetric): string => {
    const heapAllocBytes = metric?.memory?.heapAllocBytes ?? 0
    const heapSysBytes = metric?.memory?.heapSysBytes ?? 0

    return `${formatBytes(heapAllocBytes).split(' ')[0]} / ${formatBytes(heapSysBytes)}`
}

export const getHeapUsedPercent = (metric?: ExodusMetric): number => {
    const heapAllocBytes = metric?.memory?.heapAllocBytes ?? 0
    const heapSysBytes = metric?.memory?.heapSysBytes ?? 0

    if (heapSysBytes <= 0) {
        return 0
    }

    return Math.min(100, Math.max(0, (heapAllocBytes / heapSysBytes) * 100))
}

export const getExodusSummaryMetrics = (
    exodusSummary: ExodusSummary | undefined,
    exodusMetrics: ExodusMetric[] | undefined,
    t: TFunction
): IMetricCardProps[] => {
    const metric = getFirstExodusMetric(exodusMetrics)

    if (!exodusSummary && !metric) {
        return []
    }

    const uptimeSeconds = getUptimeSeconds(
        exodusSummary?.startedAt ?? metric?.startedAt,
        exodusSummary?.uptimeSeconds ?? metric?.uptimeSeconds
    )
    const totalRssBytes = exodusSummary?.totalRssBytes ?? metric?.memory?.rssBytes ?? 0
    const averageCpuPercent = exodusSummary?.averageCpuPercent ?? metric?.cpu?.processPercent ?? 0
    const schedulerP99Ms =
        exodusSummary?.averageSchedulerP99Ms ?? metric?.scheduler?.schedulerP99Ms ?? 0

    return [
        {
            value: formatCompactDuration(uptimeSeconds, t),
            IconComponent: PiClockDuotone,
            title: t('runtime-metrics.uptime'),
            iconVariant: 'soft',
            iconColor: 'blue'
        },
        {
            value: formatBytes(totalRssBytes),
            IconComponent: PiMemoryDuotone,
            title: t('runtime-metrics.total-memory'),
            iconVariant: 'soft',
            iconColor: 'cyan'
        },
        {
            value: formatPercent(averageCpuPercent),
            IconComponent: PiCpuDuotone,
            title: t('runtime-metrics.average-cpu'),
            iconVariant: 'soft',
            iconColor: 'green'
        },
        {
            value: formatMilliseconds(schedulerP99Ms),
            IconComponent: PiTimerDuotone,
            title: t('runtime-metrics.scheduler-p99'),
            iconVariant: 'soft',
            iconColor: 'orange'
        }
    ]
}

export const getExodusProcessMetrics = (
    exodusMetrics: ExodusMetric[] | undefined,
    t: TFunction
): IMetricCardProps[] => {
    const metric = getFirstExodusMetric(exodusMetrics)

    if (!metric) {
        return []
    }

    return [
        {
            value: getHeapUsedValue(metric),
            IconComponent: PiMemoryDuotone,
            title: t('runtime-metrics.heap-used'),
            iconVariant: 'soft',
            iconColor: 'cyan'
        },
        {
            value: metric.scheduler?.goroutines ?? 0,
            IconComponent: PiGearSixDuotone,
            title: t('runtime-metrics.goroutines'),
            iconVariant: 'soft',
            iconColor: 'blue'
        }
    ]
}
