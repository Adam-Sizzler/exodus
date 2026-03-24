/* eslint-disable camelcase */

import { MRT_RowSelectionState } from 'mantine-react-table'
import { devtools } from 'zustand/middleware'

import { create } from '@shared/hocs/store-wrapper'

import { IActions, IState } from './interfaces'

const initialState: IState = {
    isDrawerOpen: false,
    uuids: [],
    tableSelection: {}
}

export const useBulkUsersActionsStore = create<IActions & IState>()(
    devtools(
        (set, get) => ({
            ...initialState,
            actions: {
                setIsDrawerOpen: (isDrawerOpen: boolean) => {
                    set(() => ({ isDrawerOpen }))
                },
                setTableSelection: (
                    tableSelectionOrUpdater:
                        | ((prev: MRT_RowSelectionState) => MRT_RowSelectionState)
                        | MRT_RowSelectionState
                ) => {
                    set((state) => ({
                        tableSelection:
                            typeof tableSelectionOrUpdater === 'function'
                                ? tableSelectionOrUpdater(state.tableSelection)
                                : tableSelectionOrUpdater
                    }))
                },
                getUuidLength: () => {
                    return Object.values(get().tableSelection).filter(Boolean).length
                },
                getUuids: (): string[] => {
                    return Object.entries(get().tableSelection)
                        .filter(([, selected]) => Boolean(selected))
                        .map(([uuid]) => uuid)
                },

                getInitialState: () => {
                    return initialState
                },
                resetState: async (): Promise<void> => {
                    set({ ...initialState })
                }
            }
        }),
        {
            name: 'bulkUsersActionsStore',
            anonymousActionType: 'bulkUsersActionsStore'
        }
    )
)

export const useBulkUsersActionsStoreIsDrawerOpen = () =>
    useBulkUsersActionsStore((state) => state.isDrawerOpen)
export const useBulkUsersActionsStoreActions = () =>
    useBulkUsersActionsStore((store) => store.actions)
export const useBulkUsersActionsStoreTableSelection = () =>
    useBulkUsersActionsStore((store) => store.tableSelection)
