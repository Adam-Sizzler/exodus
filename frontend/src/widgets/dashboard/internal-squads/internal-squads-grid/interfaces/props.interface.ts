import { GetInternalSquadsCommand } from '@exodus/backend-contract'

export interface IProps {
    internalSquads: GetInternalSquadsCommand.Response['response']['internalSquads']
}
