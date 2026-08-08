import { GetInfraBillingRecordsCommand } from '@exodus/backend-contract'
import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'

dayjs.extend(utc)

export type BillingRecord = GetInfraBillingRecordsCommand.Response['response']['records'][number]

export interface MonthGroup {
    label: string
    records: BillingRecord[]
    total: number
}

export function groupRecordsByMonth(records: BillingRecord[], language: string): MonthGroup[] {
    const groups = new Map<string, BillingRecord[]>()

    for (const record of records) {
        const key = dayjs.utc(record.billedAt).format('YYYY-MM')
        const existing = groups.get(key)
        if (existing) {
            existing.push(record)
        } else {
            groups.set(key, [record])
        }
    }

    return Array.from(groups.entries()).map(([key, groupRecords]) => ({
        label: dayjs.utc(key).locale(language).format('MMMM YYYY'),
        records: groupRecords,
        total: groupRecords.reduce((sum, record) => sum + record.amount, 0)
    }))
}
