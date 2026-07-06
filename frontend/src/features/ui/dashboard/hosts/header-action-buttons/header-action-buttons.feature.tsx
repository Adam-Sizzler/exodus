import { ActionIcon, ActionIconGroup, Group, Tooltip } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { TbCards, TbPlus, TbRefresh, TbTable } from 'react-icons/tb'

import { queryClient } from '@shared/api'
import { QueryKeys, useGetHosts } from '@shared/api/hooks'
import { HelpActionIconShared } from '@shared/ui/help-drawer/help-action-icon.shared'
import { UniversalSpotlightActionIconShared } from '@shared/ui/universal-spotlight'

import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'
import { HOSTS_VIEW_MODE } from '@entities/dashboard/view-preferences-store'

interface IProps {
    setViewMode: (viewMode: HOSTS_VIEW_MODE) => void
    viewMode: HOSTS_VIEW_MODE
}

export const HeaderActionButtonsFeature = (props: IProps) => {
    const { setViewMode, viewMode } = props

    const { t } = useTranslation()

    const openModalWithData = useModalsStoreOpenWithData()

    const { isFetching } = useGetHosts()

    const handleCreate = () => {
        openModalWithData(MODALS.CREATE_HOST_MODAL, undefined)
    }

    const handleUpdate = async () => {
        await queryClient.refetchQueries({
            queryKey: QueryKeys.hosts._def
        })
    }

    return (
        <Group grow preventGrowOverflow={false} wrap="wrap">
            <HelpActionIconShared hidden={false} screen="PAGE_HOSTS" />

            <UniversalSpotlightActionIconShared />

            <ActionIconGroup>
                <Tooltip label="Toggle view mode">
                    <ActionIcon
                        color="gray"
                        onClick={() =>
                            setViewMode(
                                viewMode === HOSTS_VIEW_MODE.TABLE
                                    ? HOSTS_VIEW_MODE.CARDS
                                    : HOSTS_VIEW_MODE.TABLE
                            )
                        }
                        size="input-md"
                        variant="soft"
                    >
                        {viewMode === HOSTS_VIEW_MODE.CARDS ? (
                            <TbTable size="24px" />
                        ) : (
                            <TbCards size="24px" />
                        )}
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>

            <ActionIconGroup>
                <Tooltip label={t('common.update')} withArrow>
                    <ActionIcon
                        loading={isFetching}
                        onClick={handleUpdate}
                        size="input-md"
                        variant="soft"
                    >
                        <TbRefresh size="24px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>

            <ActionIconGroup>
                <Tooltip label={t('header-action-buttons.feature.create-new-host')} withArrow>
                    <ActionIcon color="teal" onClick={handleCreate} size="input-md" variant="soft">
                        <TbPlus size="24px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>
        </Group>
    )
}
