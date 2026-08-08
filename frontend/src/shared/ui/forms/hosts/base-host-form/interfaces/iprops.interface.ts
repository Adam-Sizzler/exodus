import { UseFormReturnType } from '@mantine/form'
import {
    CreateHostCommand,
    GetHostsTagsCommand,
    GetNodesCommand,
    GetConfigProfilesCommand,
    GetInternalSquadsCommand,
    GetSubscriptionTemplatesCommand,
    UpdateHostCommand,
    UpdateManyHostsCommand
} from '@exodus/backend-contract'

export interface IProps<
    T extends
        | CreateHostCommand.RequestBody
        | UpdateHostCommand.RequestBody
        | UpdateManyHostsCommand.RequestBody
> {
    advancedOpened: boolean
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles']
    form: UseFormReturnType<T>
    handleSubmit: () => void
    hostTags: GetHostsTagsCommand.Response['response']['tags']
    internalSquads: GetInternalSquadsCommand.Response['response']['internalSquads']
    isSubmitting: boolean
    nodes: GetNodesCommand.Response['response']
    removeRequiredFields?: boolean
    setAdvancedOpened: (value: boolean) => void
    subscriptionTemplates: GetSubscriptionTemplatesCommand.Response['response']['templates']
    hostUuid?: string
}
