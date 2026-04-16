import { GetExternalSquadsCommand } from '@exodus/backend-contract'

export interface Props {
    externalSquads: GetExternalSquadsCommand.Response['response']['externalSquads']
}
