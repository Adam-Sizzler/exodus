import { GetInternalSquadsCommand } from '@cerberus/backend-contract'

export interface Props {
    internalSquads: GetInternalSquadsCommand.Response['response']['internalSquads']
}
