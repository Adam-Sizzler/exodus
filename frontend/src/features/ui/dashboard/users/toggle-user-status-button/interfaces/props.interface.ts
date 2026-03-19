import { GetUserByUuidCommand } from '@cerberus/backend-contract'

export interface IProps {
    user: GetUserByUuidCommand.Response['response']
}
