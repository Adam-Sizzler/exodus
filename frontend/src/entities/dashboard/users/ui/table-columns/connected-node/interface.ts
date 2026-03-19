import { GetAllNodesCommand } from '@cerberus/backend-contract'

export interface IProps {
    node: GetAllNodesCommand.Response['response'][number] | undefined
}
