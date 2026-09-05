import { ActionIcon, ActionIconGroup, Group, Tooltip } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { PiPulse } from 'react-icons/pi'
import { TbPlus, TbRefresh } from 'react-icons/tb'

import { QueryKeys, useCheckSRSLists, useGetSRSLists, useSyncSRSLists } from '@shared/api/hooks'
import { queryClient } from '@shared/api/query-client'
import { UniversalSpotlightActionIconShared } from '@shared/ui/universal-spotlight'

export interface SRSListsHeaderActionButtonsProps {
    onOpenCreateModal: () => void
    srsListUuids: string[]
}

export const SRSListsHeaderActionButtonsFeature = (props: SRSListsHeaderActionButtonsProps) => {
    const { onOpenCreateModal, srsListUuids } = props
    const { t } = useTranslation()
    const tr = (key: string, defaultValue: string) => t(key, { defaultValue })

    const { isFetching, refetch: refetchSRSLists } = useGetSRSLists()
    const { mutate: checkLists, isPending: isChecking } = useCheckSRSLists({
        mutationFns: {
            onSuccess: async () => {
                await queryClient.refetchQueries({
                    queryKey: QueryKeys.srsLists.getSRSLists.queryKey
                })
            }
        }
    })
    const { mutate: syncSRSLists, isPending: isSyncing } = useSyncSRSLists()

    return (
        <Group grow preventGrowOverflow={false} wrap="wrap">
            {srsListUuids.length > 0 && <UniversalSpotlightActionIconShared />}

            <ActionIconGroup>
                <Tooltip label={tr('srs-lists.feature.check-all', 'Check all')} withArrow>
                    <ActionIcon
                        color="cyan"
                        disabled={srsListUuids.length === 0}
                        loading={isChecking}
                        onClick={() => checkLists({ variables: { uuids: srsListUuids } })}
                        size="input-md"
                        variant="soft"
                    >
                        <PiPulse size="22px" />
                    </ActionIcon>
                </Tooltip>

                <Tooltip label={tr('srs-lists.feature.sync-all', 'Sync all to disk')} withArrow>
                    <ActionIcon
                        color="blue"
                        loading={isSyncing}
                        onClick={() => syncSRSLists({})}
                        size="input-md"
                        variant="soft"
                    >
                        <TbRefresh size="22px" />
                    </ActionIcon>
                </Tooltip>

                <Tooltip label={tr('common.action.update', 'Refresh')} withArrow>
                    <ActionIcon
                        loading={isFetching}
                        onClick={() => refetchSRSLists()}
                        size="input-md"
                        variant="soft"
                    >
                        <TbRefresh size="22px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>

            <ActionIconGroup>
                <Tooltip label={tr('srs-lists.feature.add-links', 'Add SRS Lists')} withArrow>
                    <ActionIcon
                        color="teal"
                        onClick={onOpenCreateModal}
                        size="input-md"
                        variant="soft"
                    >
                        <TbPlus size="22px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>
        </Group>
    )
}
