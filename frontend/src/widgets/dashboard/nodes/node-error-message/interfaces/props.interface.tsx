import { GetNodeCommand } from '@exodus/backend-contract'

export interface IProps {
    node: GetNodeCommand.Response['response']
}
