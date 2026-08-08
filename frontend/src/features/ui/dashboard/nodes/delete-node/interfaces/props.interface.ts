import { GetNodeCommand } from '@exodus/backend-contract'

export interface IProps {
    handleClose: () => void
    node: GetNodeCommand.Response['response']
}
