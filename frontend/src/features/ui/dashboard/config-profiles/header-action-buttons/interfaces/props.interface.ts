import { GetConfigProfilesCommand } from '@cerberus/backend-contract'

export interface IProps {
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles'] | undefined
}
