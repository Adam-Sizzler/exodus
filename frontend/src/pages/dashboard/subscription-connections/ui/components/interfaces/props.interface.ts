import { GetAllNodesCommand } from '@exodus/backend-contract'

export interface IProps {
    isLoading: boolean
    nodes: GetAllNodesCommand.Response['response'] | undefined
}
