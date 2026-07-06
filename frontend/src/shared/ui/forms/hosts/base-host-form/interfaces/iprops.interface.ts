import { UseFormReturnType } from '@mantine/form'
import {
    CreateHostCommand,
    GetAllHostTagsCommand,
    GetAllNodesCommand,
    GetConfigProfilesCommand,
    GetInternalSquadsCommand,
    GetSubscriptionTemplatesCommand,
    UpdateHostCommand,
    UpdateManyHostsCommand
} from '@exodus/backend-contract'

export interface IProps<
    T extends CreateHostCommand.Request | UpdateHostCommand.Request | UpdateManyHostsCommand.Request
> {
    advancedOpened: boolean
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles']
    form: UseFormReturnType<T>
    handleSubmit: () => void
    hostTags: GetAllHostTagsCommand.Response['response']['tags']
    internalSquads: GetInternalSquadsCommand.Response['response']['internalSquads']
    isSubmitting: boolean
    nodes: GetAllNodesCommand.Response['response']
    removeRequiredFields?: boolean
    setAdvancedOpened: (value: boolean) => void
    subscriptionTemplates: GetSubscriptionTemplatesCommand.Response['response']['templates']
}
