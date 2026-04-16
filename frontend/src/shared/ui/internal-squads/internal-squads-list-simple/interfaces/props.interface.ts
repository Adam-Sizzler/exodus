import { GetInternalSquadsCommand } from '@exodus/backend-contract'

export interface IProps {
    filteredInternalSquads: GetInternalSquadsCommand.Response['response']['internalSquads']
}
