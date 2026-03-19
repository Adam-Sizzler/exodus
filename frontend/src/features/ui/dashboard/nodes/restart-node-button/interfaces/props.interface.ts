import { GetOneNodeCommand } from '@cerberus/backend-contract'

export interface IProps {
    handleClose: () => void
    node: GetOneNodeCommand.Response['response']
}
