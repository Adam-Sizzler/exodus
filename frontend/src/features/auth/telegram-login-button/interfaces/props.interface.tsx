import { GetStatusCommand } from '@cerberus/backend-contract'

export interface IProps {
    authentication: NonNullable<GetStatusCommand.Response['response']['authentication']>
}
