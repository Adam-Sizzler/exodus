import { GetConfigProfileByUuidCommand } from '@exodus/backend-contract'

export interface IProps {
    configProfile: GetConfigProfileByUuidCommand.Response['response']
}
