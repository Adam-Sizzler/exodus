import { GetAllNodesCommand } from '@exodus/backend-contract'

export interface IProps {
    nodes: GetAllNodesCommand.Response['response'] | undefined
}
