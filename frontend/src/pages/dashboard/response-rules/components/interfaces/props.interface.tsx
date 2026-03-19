import { GetConfigProfilesCommand } from '@cerberus/backend-contract'

export interface Props {
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles']
}
