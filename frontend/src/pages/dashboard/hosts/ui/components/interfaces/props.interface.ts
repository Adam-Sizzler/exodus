import {
    GetHostsCommand,
    GetHostsTagsCommand,
    GetConfigProfilesCommand
} from '@exodus/backend-contract'

export interface IProps {
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles'] | undefined
    hosts: GetHostsCommand.Response['response'] | undefined
    hostTags: GetHostsTagsCommand.Response['response']['tags'] | undefined
    isLoading: boolean
}
