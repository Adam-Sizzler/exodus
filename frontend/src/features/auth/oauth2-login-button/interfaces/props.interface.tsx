import { GetStatusCommand } from '@exodus/backend-contract'

export interface IProps {
    authentication: NonNullable<GetStatusCommand.Response['response']['authentication']>
}
