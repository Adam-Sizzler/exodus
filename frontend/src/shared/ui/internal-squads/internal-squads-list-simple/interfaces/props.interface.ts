import { GetInternalSquadsCommand } from '@cerberus/backend-contract'

export interface IProps {
    filteredInternalSquads: GetInternalSquadsCommand.Response['response']['internalSquads']
}
