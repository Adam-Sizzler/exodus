import type { IProps } from './interfaces/props.interface'

import { Accordion, Badge, Divider, Group, Text, Tooltip } from '@mantine/core'
import { memo } from 'react'
import { useTranslation } from 'react-i18next'
import { PiCheckBold, PiCpu } from 'react-icons/pi'

import { showModal } from '@shared/_modals/show-modal'

import { VirtualizedInboundsListShared } from '../virtualized-inbounds-list/virtualized-inbounds-list.shared'
import { AccordionControlShared } from './accordion-control.shared'

export const ConfigProfileCardShared = memo((props: IProps) => {
    const {
        hideSelectActions,
        profile,
        selectedInbounds,
        onInboundToggle,
        onSelectAllInbounds,
        onUnselectAllInbounds,
        isOpen
    } = props

    const { t } = useTranslation()

    const selectedInboundsFromProfile = profile.inbounds.filter((inbound) =>
        selectedInbounds.has(inbound.uuid)
    ).length

    return (
        <Accordion.Item
            style={{
                border: '1px solid rgb(255, 255, 255, 0.08)',
                background: 'rgb(255, 255, 255, 0.02)'
            }}
            value={profile.uuid}
        >
            <AccordionControlShared
                hideSelectActions={hideSelectActions}
                onSelectAllInbounds={onSelectAllInbounds}
                onUnselectAllInbounds={onUnselectAllInbounds}
                profile={profile}
                value={profile.uuid}
            >
                <Group mb="xs">
                    <Text ff="monospace" fw={700} size="md">
                        {profile.name}
                    </Text>
                </Group>
                <Group>
                    <Badge
                        color={
                            selectedInboundsFromProfile === profile.inbounds.length
                                ? 'teal'
                                : 'cyan'
                        }
                        ff="monospace"
                        leftSection={<PiCheckBold />}
                        size="md"
                        variant="soft"
                    >
                        {selectedInboundsFromProfile} / {profile.inbounds.length}
                    </Badge>

                    <Tooltip label={t('config-profile-card.shared.active-on-nodes')}>
                        <Badge
                            color={profile.nodes.length > 0 ? 'teal' : 'cyan'}
                            ff="monospace"
                            leftSection={<PiCpu />}
                            onClick={(e) => {
                                e.preventDefault()
                                e.stopPropagation()
                                e.nativeEvent.stopImmediatePropagation()

                                showModal('configProfiles_activeNodesModal', {
                                    nodes: profile.nodes,
                                    profileName: profile.name
                                })
                            }}
                            size="md"
                            style={{ cursor: 'pointer' }}
                            variant="soft"
                        >
                            {profile.nodes.length}
                        </Badge>
                    </Tooltip>
                </Group>
            </AccordionControlShared>
            <Accordion.Panel>
                <Divider p="xs" style={{ opacity: 0.3 }} />
                {isOpen && (
                    <VirtualizedInboundsListShared
                        onInboundToggle={onInboundToggle}
                        profile={profile}
                        selectedInbounds={selectedInbounds}
                    />
                )}
            </Accordion.Panel>
        </Accordion.Item>
    )
})
