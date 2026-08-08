import { GetNodesCommand } from '@exodus/backend-contract'

export interface IProps {
    node: GetNodesCommand.Response['response'][number] | undefined
}
