import { MRT_RowSelectionState } from '@kastov/mantine-react-table-open'
import { devtools } from 'zustand/middleware'

import { create } from '@shared/hocs/store-wrapper'

import { IActions, IState } from './interfaces'

const initialState: IState = {
    userIds: [],
    tableSelection: {}
}

export const useUsersTableSelectionStore = create<IActions & IState>()(
    devtools(
        (set, get) => ({
            ...initialState,
            actions: {
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
                getIdsLength: () => {
                    return Object.keys(get().tableSelection).length
                },
                getIds: (): number[] => {
                    return Object.keys(get().tableSelection).map((id) => Number(id))
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
            name: 'usersTableSelectionStore',
            anonymousActionType: 'usersTableSelectionStore'
        }
    )
)

export const useUsersTableSelectionStoreActions = () =>
    useUsersTableSelectionStore((store) => store.actions)
export const useUsersTableSelectionStoreTableSelection = () =>
    useUsersTableSelectionStore((store) => store.tableSelection)
