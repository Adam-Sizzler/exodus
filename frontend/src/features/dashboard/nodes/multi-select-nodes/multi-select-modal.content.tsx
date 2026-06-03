import { NODES_BULK_ACTIONS, TNodesBulkActions } from '@exodus/backend-contract'
import { TbCheck, TbPlayerStop, TbRefresh, TbRocket } from 'react-icons/tb'
import { useTranslation } from 'react-i18next'
import { modals } from '@mantine/modals'
import { Stack } from '@mantine/core'

import { NodeResponse, QueryKeys, useBulkNodesActions } from '@shared/api/hooks'
import { queryClient } from '@shared/api/query-client'
import { ActionCardShared } from '@shared/ui'

type NodeType = NodeResponse

interface IProps {
    selectedRecords: NodeType[]
    setSelectedRecords: (records: NodeType[]) => void
}

export const MultiSelectNodesModalContent = (props: IProps) => {
    const { selectedRecords, setSelectedRecords } = props
    const { t } = useTranslation()
    const { mutateAsync: bulkAction, isPending } = useBulkNodesActions()

    const uuids = selectedRecords.map((node) => node.uuid)

    const handleAction = async (action: TNodesBulkActions) => {
        if (isPending || uuids.length === 0) return
        await bulkAction({ variables: { uuids, action } })

        queryClient.refetchQueries({ queryKey: QueryKeys.nodes.getAllNodes.queryKey })
        modals.closeAll()
        setSelectedRecords([])
    }

    return (
        <Stack gap="xs">
            <ActionCardShared
                description={`${uuids.length} node(s)`}
                icon={<TbRocket size={20} />}
                isLoading={isPending}
                onClick={() => handleAction(NODES_BULK_ACTIONS.RESTART)}
                title={t('restart-node-button.feature.restart')}
                iconColor="teal"
                variant="soft"
            />
            <ActionCardShared
                description={`${uuids.length} node(s)`}
                icon={<TbPlayerStop size={20} />}
                isLoading={isPending}
                onClick={() => handleAction(NODES_BULK_ACTIONS.DISABLE)}
                title={t('common.disable')}
                iconColor="orange"
                variant="soft"
            />
            <ActionCardShared
                description={`${uuids.length} node(s)`}
                icon={<TbCheck size={20} />}
                isLoading={isPending}
                onClick={() => handleAction(NODES_BULK_ACTIONS.ENABLE)}
                title={t('common.enable')}
                iconColor="cyan"
                variant="soft"
            />
            <ActionCardShared
                description={`${uuids.length} node(s)`}
                icon={<TbRefresh size={20} />}
                isLoading={isPending}
                onClick={() => handleAction(NODES_BULK_ACTIONS.RESET_TRAFFIC)}
                title={t('reset-node-traffic.feature.reset-traffic')}
                iconColor="violet"
                variant="soft"
            />
        </Stack>
    )
}
