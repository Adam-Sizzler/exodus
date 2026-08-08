import HostsPageComponent from '@pages/dashboard/hosts/ui/components/hosts.page.component'

import {
    useGetConfigProfiles,
    useGetHosts,
    useGetHostTags,
    useGetInternalSquads,
    useGetNodes,
    useGetSubscriptionTemplates
} from '@shared/api/hooks'

export function HostsPageConnector() {
    const { data: hosts, isLoading: isHostsLoading } = useGetHosts()
    const { isLoading: isInternalSquadsLoading } = useGetInternalSquads()
    const { data: configProfiles, isLoading: isConfigProfilesLoading } = useGetConfigProfiles()
    const { data: hostTags, isLoading: isHostTagsLoading } = useGetHostTags()
    const { isLoading: isNodesLoading } = useGetNodes()
    const { isLoading: isSubscriptionTemplatesLoading } = useGetSubscriptionTemplates()

    return (
        <HostsPageComponent
            configProfiles={configProfiles?.configProfiles}
            hosts={hosts}
            hostTags={hostTags?.tags}
            isLoading={
                isConfigProfilesLoading ||
                isHostsLoading ||
                isHostTagsLoading ||
                isNodesLoading ||
                isSubscriptionTemplatesLoading ||
                isInternalSquadsLoading
            }
        />
    )
}
