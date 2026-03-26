import { GetAllNodesCommand } from '@cerberus/backend-contract'

export interface IProps {
    nodes: GetAllNodesCommand.Response['response'] | undefined
}
