import { GetInternalSquadsCommand } from '@cerberus/backend-contract'

export interface IProps {
    hideEditButton?: boolean
    internalSquad: GetInternalSquadsCommand.Response['response']['internalSquads'][number] | null
}
