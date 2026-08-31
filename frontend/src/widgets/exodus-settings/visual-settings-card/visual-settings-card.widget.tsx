import { Group, SegmentedControl, Stack } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { TbLayoutNavbar, TbLayoutSidebar } from 'react-icons/tb'

import { useIsMobile } from '@shared/hooks'
import { SettingsCardShared } from '@shared/ui/settings-card'

import {
    useExperimentalFeature,
    useViewPreferencesStoreActions
} from '@entities/dashboard/view-preferences-store'

export const VisualSettingsCardWidget = () => {
    const { t } = useTranslation()
    const isMobile = useIsMobile()

    const isLegacyLayoutStyle = useExperimentalFeature('legacyLayoutStyle')
    const { setExperimentalFeature } = useViewPreferencesStoreActions()

    if (isMobile) {
        return null
    }

    return (
        <SettingsCardShared.Container>
            <SettingsCardShared.Header
                description={t(
                    'visual-settings-card.widget.choose-how-the-desktop-navigation-is-displayed'
                )}
                icon={<TbLayoutSidebar size={24} />}
                iconColor="cyan"
                iconVariant="soft"
                title={t('visual-settings-card.widget.layout-style')}
            />

            <SettingsCardShared.Content>
                <Stack gap="md">
                    <SegmentedControl
                        data={[
                            {
                                label: (
                                    <Group gap="xs" justify="center" wrap="nowrap">
                                        <TbLayoutNavbar size={18} />
                                        <span>{t('visual-settings-card.widget.compact')}</span>
                                    </Group>
                                ),
                                value: 'compact'
                            },
                            {
                                label: (
                                    <Group gap="xs" justify="center" wrap="nowrap">
                                        <TbLayoutSidebar size={18} />
                                        <span>{t('visual-settings-card.widget.sidebar')}</span>
                                    </Group>
                                ),
                                value: 'sidebar'
                            }
                        ]}
                        fullWidth
                        onChange={(value) =>
                            setExperimentalFeature('legacyLayoutStyle', value === 'sidebar')
                        }
                        value={isLegacyLayoutStyle ? 'sidebar' : 'compact'}
                    />
                </Stack>
            </SettingsCardShared.Content>
        </SettingsCardShared.Container>
    )
}
