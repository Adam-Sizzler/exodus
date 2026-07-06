/* eslint-disable camelcase */
import {
    MRT_ColumnFilterFnsState,
    MRT_ColumnFiltersState,
    MRT_ColumnOrderState,
    MRT_ColumnPinningState,
    MRT_SortingState,
    MRT_VisibilityState
} from '@kastov/mantine-react-table-open'
import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'

export interface UsersTableTemplateSnapshot {
    columnFilterFns: MRT_ColumnFilterFnsState
    columnFilters: MRT_ColumnFiltersState
    columnOrder: MRT_ColumnOrderState
    columnPinning: MRT_ColumnPinningState
    columnVisibility: MRT_VisibilityState
    showColumnFilters: boolean
    sorting: MRT_SortingState
}

export interface UsersTableTemplate {
    color: string
    id: string
    name: string
    snapshot: UsersTableTemplateSnapshot
}

export interface UsersTableTemplatesActions {
    addTemplate: (name: string, color: string, snapshot: UsersTableTemplateSnapshot) => void
    removeTemplate: (id: string) => void
}

export interface UsersTableTemplatesStore {
    actions: UsersTableTemplatesActions
    templates: UsersTableTemplate[]
}

export const useUsersTableTemplatesStore = create<UsersTableTemplatesStore>()(
    persist(
        (set) => ({
            templates: [],
            actions: {
                addTemplate: (name, color, snapshot) =>
                    set((state) => ({
                        templates: [
                            ...state.templates,
                            { id: crypto.randomUUID(), name, color, snapshot }
                        ]
                    })),
                removeTemplate: (id) =>
                    set((state) => ({
                        templates: state.templates.filter((template) => template.id !== id)
                    }))
            }
        }),
        {
            name: 'x-rmnw-users-table-templates',
            version: 1,
            storage: createJSONStorage(() => localStorage),
            partialize: (state) => ({ templates: state.templates })
        }
    )
)

export const useUsersTableTemplates = () => useUsersTableTemplatesStore((store) => store.templates)

export const useUsersTableTemplatesActions = () =>
    useUsersTableTemplatesStore((store) => store.actions)
