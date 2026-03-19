import { GetExternalSquadsCommand } from '@cerberus/backend-contract'

export interface IProps {
    externalSquads: GetExternalSquadsCommand.Response['response']['externalSquads']
}
