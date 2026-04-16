import { GetOneNodeCommand } from '@exodus/backend-contract'

export interface IProps {
    node: GetOneNodeCommand.Response['response']
}
