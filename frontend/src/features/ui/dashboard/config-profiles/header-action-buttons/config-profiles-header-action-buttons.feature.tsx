import { ActionIcon, ActionIconGroup, Group, Tooltip } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { TbCode, TbPlus, TbRefresh } from 'react-icons/tb'
import { useNavigate } from 'react-router'

import { showModal } from '@shared/_modals/show-modal'
import { HelpActionIconShared } from '@shared/_modals/universal'
import { queryClient } from '@shared/api'
import { QueryKeys, useGetConfigProfiles } from '@shared/api/hooks'
import { SingboxLogo } from '@shared/ui/logos'
import { UniversalSpotlightActionIconShared } from '@shared/ui/universal-spotlight'

import { CONFIG_PROFILES_VIEW_MODE } from '@entities/dashboard/view-preferences-store'

interface IProps {
    configProfileCount: number
    setViewMode: (viewMode: CONFIG_PROFILES_VIEW_MODE) => void
    viewMode: CONFIG_PROFILES_VIEW_MODE
}

export const ConfigProfilesHeaderActionButtonsFeature = (props: IProps) => {
    const { configProfileCount, setViewMode, viewMode } = props
    const { isFetching } = useGetConfigProfiles()
    const { t } = useTranslation()

    const navigate = useNavigate()

    const handleUpdate = async () => {
        await queryClient.refetchQueries({
            queryKey: QueryKeys.configProfiles.getConfigProfiles.queryKey
        })
    }

    return (
        <Group grow preventGrowOverflow={false} wrap="wrap">
            <HelpActionIconShared hidden={false} screen="PAGE_CONFIG_PROFILES" />

            {configProfileCount > 0 && <UniversalSpotlightActionIconShared />}

            <ActionIconGroup>
                <ActionIcon
                    color="gray"
                    onClick={() =>
                        setViewMode(
                            viewMode === CONFIG_PROFILES_VIEW_MODE.PROFILES
                                ? CONFIG_PROFILES_VIEW_MODE.SNIPPETS
                                : CONFIG_PROFILES_VIEW_MODE.PROFILES
                        )
                    }
                    size="input-md"
                    variant="soft"
                >
                    {viewMode === CONFIG_PROFILES_VIEW_MODE.PROFILES ? (
                        <TbCode size="24px" />
                    ) : (
                        <SingboxLogo size="24px" />
                    )}
                </ActionIcon>
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
                <Tooltip
                    label={t('config-profiles-header-action-buttons.feature.create-config-profile')}
                    withArrow
                >
                    <ActionIcon
                        color="teal"
                        onClick={() =>
                            showModal('createModal', {
                                createFrom: 'configProfile',
                                contentOptions: { navigate }
                            })
                        }
                        size="input-md"
                        variant="soft"
                    >
                        <TbPlus size="24px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>
        </Group>
    )
}
