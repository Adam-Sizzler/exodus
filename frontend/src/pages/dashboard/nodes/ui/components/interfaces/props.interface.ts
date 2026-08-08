import { GetNodesCommand } from '@exodus/backend-contract'

export interface IProps {
    isLoading: boolean
    nodes: GetNodesCommand.Response['response'] | undefined
}
