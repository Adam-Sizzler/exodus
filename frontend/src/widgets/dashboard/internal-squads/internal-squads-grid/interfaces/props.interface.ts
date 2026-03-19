import { GetInternalSquadsCommand } from '@cerberus/backend-contract'

export interface IProps {
    internalSquads: GetInternalSquadsCommand.Response['response']['internalSquads']
}
