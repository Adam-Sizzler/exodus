import {
    PiClockDuotone,
    PiCloudDuotone,
    PiCpuDuotone,
    PiGearSixDuotone,
    PiMemoryDuotone,
    PiQueueDuotone,
    PiRocketLaunchDuotone
} from 'react-icons/pi'
import { GetExodusHealthCommand } from '@exodus/backend-contract'
import { ThemeIconProps } from '@mantine/core'
import { TFunction } from 'i18next'

import { IMetricCardProps } from '@shared/ui/metrics/metric-card'
import { prettyBytesToAnyUtil } from '@shared/utils/bytes'

const getProcessIcon = (processName: string) => {
    if (processName.includes('REST API')) return PiCloudDuotone
    if (processName.includes('Scheduler')) return PiClockDuotone
    if (processName.includes('Jobs')) return PiQueueDuotone
    return PiGearSixDuotone
}

const getIconColor = (cpuUsage: number): ThemeIconProps['color'] => {
    if (cpuUsage < 30) return 'green'
    if (cpuUsage < 70) return 'yellow'
    return 'red'
}

export const getPm2SummaryMetrics = (
    pm2Stats: GetExodusHealthCommand.Response['response']['pm2Stats'],
    t: TFunction
): IMetricCardProps[] => {
    if (!pm2Stats || pm2Stats.length === 0) {
        return []
    }

    type PM2Process = GetExodusHealthCommand.Response['response']['pm2Stats'][number]

    const totalMemoryBytes = pm2Stats.reduce((sum: number, process: PM2Process) => {
        return sum + Number(process.memory)
    }, 0)

    const averageCpu =
        pm2Stats.reduce((sum: number, process: PM2Process) => {
            return sum + parseFloat(process.cpu)
        }, 0) / pm2Stats.length

    const heaviestProcess = pm2Stats.reduce((max: PM2Process, process: PM2Process) => {
        return parseFloat(process.cpu) > parseFloat(max.cpu) ? process : max
    })

    return [
        {
            value: pm2Stats.length,
            IconComponent: PiGearSixDuotone,
            title: t('pm2-metrics.total-processes'),
            iconVariant: 'soft',
            iconColor: 'blue'
        },
        {
            value: prettyBytesToAnyUtil(totalMemoryBytes, true),
            IconComponent: PiMemoryDuotone,
            title: t('pm2-metrics.total-memory'),
            iconVariant: 'soft',
            iconColor: 'cyan'
        },
        {
            value: `${averageCpu.toFixed(1)}%`,
            IconComponent: PiCpuDuotone,
            title: t('pm2-metrics.average-cpu'),
            iconVariant: 'soft',
            iconColor: 'green'
        },
        {
            value: `${heaviestProcess.name}`,
            IconComponent: PiRocketLaunchDuotone,
            title: t('pm2-metrics.heaviest-process'),
            iconVariant: 'soft',
            iconColor: 'orange'
        }
    ]
}

export const getPm2ProcessMetrics = (
    pm2Stats: GetExodusHealthCommand.Response['response']['pm2Stats']
): IMetricCardProps[] => {
    if (!pm2Stats || pm2Stats.length === 0) {
        return []
    }

    return pm2Stats.map((process: GetExodusHealthCommand.Response['response']['pm2Stats'][number]) => {
        const cpuUsage = parseFloat(process.cpu)
        return {
            value: prettyBytesToAnyUtil(process.memory, true),
            IconComponent: getProcessIcon(process.name),
            title: process.name,
            iconVariant: 'soft',
            iconColor: getIconColor(cpuUsage)
        }
    })
}

export const getPm2Metrics = getPm2SummaryMetrics
