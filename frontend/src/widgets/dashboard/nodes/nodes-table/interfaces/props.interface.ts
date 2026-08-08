import { GetNodesCommand } from '@exodus/backend-contract'

export interface IProps {
    nodes: GetNodesCommand.Response['response'] | undefined
}
