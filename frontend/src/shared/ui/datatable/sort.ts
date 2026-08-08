import { type DataTableSortStatus } from '@kastov/mantine-datatable'
import get from 'lodash/get'

type SortStatus<T> = Pick<DataTableSortStatus<T>, 'columnAccessor' | 'direction'>

export function sortRecords<T>(records: T[], sortStatus: SortStatus<T>): T[] {
    const isDesc = sortStatus.direction === 'desc'

    return [...records].sort((a, b) => {
        const aVal = get(a, sortStatus.columnAccessor)
        const bVal = get(b, sortStatus.columnAccessor)

        if (aVal == null && bVal == null) return 0
        if (aVal == null) return 1
        if (bVal == null) return -1

        let result: number
        if (typeof aVal === 'string' && typeof bVal === 'string') {
            result = aVal.toLowerCase().localeCompare(bVal.toLowerCase())
        } else if (aVal < bVal) {
            result = -1
        } else {
            result = aVal > bVal ? 1 : 0
        }
        return isDesc ? -result : result
    })
}
