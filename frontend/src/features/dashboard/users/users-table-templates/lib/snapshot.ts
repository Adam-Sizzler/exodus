/* eslint-disable camelcase */
import { MRT_TableInstance } from '@kastov/mantine-react-table-open'
import { GetAllUsersCommand } from '@exodus/backend-contract'

import { UsersTableTemplateSnapshot } from '@entities/dashboard/users/users-table-templates-store'

type UserRow = GetAllUsersCommand.Response['response']['users'][number]
type UsersTable = MRT_TableInstance<UserRow>

export const captureSnapshot = (table: UsersTable): UsersTableTemplateSnapshot => {
    const state = table.getState()

    return {
        columnFilterFns: { ...state.columnFilterFns },
        columnFilters: [...state.columnFilters],
        columnOrder: [...state.columnOrder],
        columnPinning: {
            left: [...(state.columnPinning.left ?? [])],
            right: [...(state.columnPinning.right ?? [])]
        },
        columnVisibility: { ...state.columnVisibility },
        showColumnFilters: state.showColumnFilters,
        sorting: [...state.sorting]
    }
}

export const applySnapshot = (table: UsersTable, snapshot: UsersTableTemplateSnapshot) => {
    table.setShowColumnFilters(snapshot.showColumnFilters)
    table.setColumnVisibility(snapshot.columnVisibility)
    table.setColumnOrder(snapshot.columnOrder)
    table.setColumnPinning(snapshot.columnPinning)
    table.setColumnFilterFns(snapshot.columnFilterFns)
    table.setColumnFilters(snapshot.columnFilters)
    table.setSorting(snapshot.sorting)
    table.setPageIndex(0)
}
