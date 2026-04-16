import { GetConfigProfilesCommand } from '@exodus/backend-contract'

export interface Props {
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles']
}
