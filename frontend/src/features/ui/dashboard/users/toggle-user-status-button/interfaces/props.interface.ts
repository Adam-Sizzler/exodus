import { GetUserByUuidCommand } from '@exodus/backend-contract'

export interface IProps {
    user: GetUserByUuidCommand.Response['response']
}
