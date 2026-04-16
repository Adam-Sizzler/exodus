import { GetAllNodesCommand } from '@exodus/backend-contract'

export interface IProps {
    node: GetAllNodesCommand.Response['response'][number] | undefined
}
