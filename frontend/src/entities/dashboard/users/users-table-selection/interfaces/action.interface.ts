import { MRT_RowSelectionState } from '@kastov/mantine-react-table-open'

export interface IActions {
    actions: {
        getIdsLength: () => number
        getIds: () => number[]
        resetState: () => void
        setTableSelection: (
            tableSelectionOrUpdater:
                | ((prev: MRT_RowSelectionState) => MRT_RowSelectionState)
                | MRT_RowSelectionState
        ) => void
    }
}
