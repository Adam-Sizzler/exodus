import { GetExternalSquadsCommand } from '@cerberus/backend-contract'

export interface Props {
    externalSquads: GetExternalSquadsCommand.Response['response']['externalSquads']
}
