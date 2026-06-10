import { PiClockDuotone, PiCpuDuotone, PiGearSixDuotone, PiMemoryDuotone } from 'react-icons/pi'
import { TFunction } from 'i18next'

import { IMetricCardProps } from '@shared/ui/metrics/metric-card'
import { prettyBytesToAnyUtil } from '@shared/utils/bytes'

export type RuntimeSummary = {
    uptimeSeconds?: number
    totalRssBytes?: number
    heapAllocBytes?: number
    averageCpuPercent?: number
    averageGoroutines?: number
    averageSchedulerDelayMs?: number
    averageSchedulerP99Ms?: number
}

export type RuntimeMetric = {
    name?: string
    instanceType?: string
    instanceId?: number
    pid?: number
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

export const getFirstRuntimeMetric = (runtimeMetrics?: RuntimeMetric[]): RuntimeMetric | undefined => {
    if (!runtimeMetrics || runtimeMetrics.length === 0) {
        return undefined
    }

    return runtimeMetrics[0]
}

export const formatMilliseconds = (value?: number): string => `${Number(value ?? 0).toFixed(3)} ms`
const formatPercent = (value?: number): string => `${Number(value ?? 0).toFixed(1)}%`
export const formatBytes = (value?: number): string => prettyBytesToAnyUtil(Number(value ?? 0), true)

const formatCompactDuration = (uptimeInSeconds?: number): string => {
    const totalSeconds = Math.max(0, Math.floor(Number(uptimeInSeconds ?? 0)))
    const days = Math.floor(totalSeconds / 86_400)
    const hours = Math.floor((totalSeconds % 86_400) / 3_600)
    const minutes = Math.floor((totalSeconds % 3_600) / 60)
    const seconds = totalSeconds % 60

    const parts: string[] = []

    if (days > 0) {
        parts.push(`${days}d`)
    }
    if (hours > 0) {
        parts.push(`${hours}h`)
    }
    if (minutes > 0) {
        parts.push(`${minutes}m`)
    }
    if (seconds > 0 || parts.length === 0) {
        parts.push(`${seconds}s`)
    }

    return parts.join(' ')
}

export const getHeapUsedValue = (metric?: RuntimeMetric): string => {
    const heapAllocBytes = metric?.memory?.heapAllocBytes ?? 0
    const heapSysBytes = metric?.memory?.heapSysBytes ?? 0

    return `${formatBytes(heapAllocBytes)} / ${formatBytes(heapSysBytes)}`
}

export const getHeapUsedPercent = (metric?: RuntimeMetric): number => {
    const heapAllocBytes = metric?.memory?.heapAllocBytes ?? 0
    const heapSysBytes = metric?.memory?.heapSysBytes ?? 0

    if (heapSysBytes <= 0) {
        return 0
    }

    return Math.min(100, Math.max(0, (heapAllocBytes / heapSysBytes) * 100))
}

export const getRuntimeSummaryMetrics = (
    runtimeSummary: RuntimeSummary | undefined,
    runtimeMetrics: RuntimeMetric[] | undefined,
    t: TFunction
): IMetricCardProps[] => {
    const metric = getFirstRuntimeMetric(runtimeMetrics)

    if (!runtimeSummary && !metric) {
        return []
    }

    const uptimeSeconds = runtimeSummary?.uptimeSeconds ?? metric?.uptimeSeconds ?? 0
    const totalRssBytes = runtimeSummary?.totalRssBytes ?? metric?.memory?.rssBytes ?? 0
    const averageCpuPercent = runtimeSummary?.averageCpuPercent ?? metric?.cpu?.processPercent ?? 0
    const schedulerP99Ms = runtimeSummary?.averageSchedulerP99Ms ?? metric?.scheduler?.schedulerP99Ms ?? 0

    return [
        {
            value: formatCompactDuration(uptimeSeconds),
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
            IconComponent: PiClockDuotone,
            title: t('runtime-metrics.scheduler-p99'),
            iconVariant: 'soft',
            iconColor: 'orange'
        }
    ]
}

export const getRuntimeProcessMetrics = (
    runtimeMetrics: RuntimeMetric[] | undefined,
    t: TFunction
): IMetricCardProps[] => {
    const metric = getFirstRuntimeMetric(runtimeMetrics)

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

// Backward-compatible aliases for older imports. New UI should use runtimeSummary/runtimeMetrics.
export const getPm2SummaryMetrics = (..._args: unknown[]): IMetricCardProps[] => []
export const getPm2ProcessMetrics = (..._args: unknown[]): IMetricCardProps[] => []
export const getPm2Metrics = getPm2SummaryMetrics
