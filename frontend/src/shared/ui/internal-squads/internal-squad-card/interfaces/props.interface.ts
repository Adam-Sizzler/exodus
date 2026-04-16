import { GetInternalSquadsCommand } from '@exodus/backend-contract'

export interface IProps {
    internalSquad: GetInternalSquadsCommand.Response['response']['internalSquads'][number] | null
}
