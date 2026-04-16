import { GetInternalSquadsCommand } from '@exodus/backend-contract'

export interface Props {
    internalSquads: GetInternalSquadsCommand.Response['response']['internalSquads']
}
