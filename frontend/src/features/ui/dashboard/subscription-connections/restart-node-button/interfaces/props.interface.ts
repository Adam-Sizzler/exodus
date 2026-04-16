import { GetOneNodeCommand } from '@exodus/backend-contract'

export interface IProps {
    handleClose: () => void
    node: GetOneNodeCommand.Response['response']
}
