import { GetAllNodesCommand } from '@cerberus/backend-contract'

export interface IProps {
    isLoading: boolean
    nodes: GetAllNodesCommand.Response['response'] | undefined
}
