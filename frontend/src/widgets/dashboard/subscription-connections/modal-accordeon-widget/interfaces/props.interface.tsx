import { GetOneNodeCommand } from '@cerberus/backend-contract'

export interface IProps {
    node: GetOneNodeCommand.Response['response']
}
