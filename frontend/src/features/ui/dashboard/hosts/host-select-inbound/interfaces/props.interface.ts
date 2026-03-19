import { GetConfigProfilesCommand } from '@cerberus/backend-contract'

export interface IProps {
    activeConfigProfileInbound: null | string | undefined
    activeConfigProfileUuid: null | string | undefined
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles']
    onSaveInbound: (inbound: string, configProfileUuid: string) => void
}
