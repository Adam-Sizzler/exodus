import { useEffect } from 'react'

import {
    subscriptionConnectionsQueryKeys,
    QueryKeys,
    useGetConfigProfiles,
    useGetSubNodes,
    useGetSubscriptionConnectionsTags
} from '@shared/api/hooks'
import {
    useSubscriptionConnectionsStoreActions,
    useSubscriptionConnectionsStoreCreateModalIsOpen
} from '@entities/dashboard/subscription-connections/nodes-store'
import { queryClient } from '@shared/api'

import NodesPageComponent from '../components/nodes.page.component'

export function SubscriptionConnectionsPageConnector() {
    const actions = useSubscriptionConnectionsStoreActions()

    const isCreateModalOpen = useSubscriptionConnectionsStoreCreateModalIsOpen()

    const { data: nodes, isLoading } = useGetSubNodes()
    const { isLoading: isConfigProfilesLoading } = useGetConfigProfiles()
    useGetSubscriptionConnectionsTags()

    useEffect(() => {
        ;(async () => {
            await queryClient.prefetchQuery({
                queryKey: subscriptionConnectionsQueryKeys.getPubKey.queryKey
            })
        })()
        return () => {
            actions.resetState()
        }
    }, [])

    useEffect(() => {
        if (isCreateModalOpen) return
        ;(async () => {
            await queryClient.refetchQueries({
                queryKey: QueryKeys.subscriptionConnections.getAllNodes.queryKey
            })
        })()
    }, [isCreateModalOpen])

    return <NodesPageComponent isLoading={isLoading || isConfigProfilesLoading} nodes={nodes} />
}
