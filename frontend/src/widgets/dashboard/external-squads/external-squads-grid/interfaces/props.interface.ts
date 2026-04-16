import { GetExternalSquadsCommand } from '@exodus/backend-contract'

export interface IProps {
    externalSquads: GetExternalSquadsCommand.Response['response']['externalSquads']
}
