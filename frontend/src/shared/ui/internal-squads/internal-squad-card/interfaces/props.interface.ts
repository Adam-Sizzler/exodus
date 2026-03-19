import { GetInternalSquadsCommand } from '@cerberus/backend-contract'

export interface IProps {
    internalSquad: GetInternalSquadsCommand.Response['response']['internalSquads'][number] | null
}
