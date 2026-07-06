import { Box, Center, DataList, Drawer, Group, Stack } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import {
    PiArrowsDownUpDuotone,
    PiCalendarDotDuotone,
    PiClockDuotone,
    PiNetworkDuotone,
    PiTagDuotone,
    PiUserDuotone
} from 'react-icons/pi'

import { useGetUserByUuid } from '@shared/api/hooks'
import { useGetSubscriptionConnections } from '@shared/api/hooks/subscription-connections/subscription-connections.query.hooks'
import { CopyableDataListItem } from '@shared/ui/copyable-field/copyable-data-list-item'
import { LoaderModalShared } from '@shared/ui/loader-modal'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCardRoot } from '@shared/ui/section-card/section-card.root'
import { SectionCardSection } from '@shared/ui/section-card/section-card.section'
import { prettifyBytesUtil } from '@shared/utils/bytes'
import { buildSubscriptionLinksFromNodes } from '@shared/utils/misc/subscription-links'
import { formatTimeUtil } from '@shared/utils/time-utils'

import {
    useUserModalStoreActions,
    useUserModalStoreDrawerUserUuid,
    useUserModalStoreIsDetailedUserInfoDrawerOpen
} from '@entities/dashboard/user-modal-store/user-modal-store'

import { UserStatusBadge } from '../user-status-badge/user-status-badge.widget'

export const DetailedUserInfoDrawerWidget = () => {
    const { t, i18n } = useTranslation()

    const actions = useUserModalStoreActions()
    const isDetailedUserInfoDrawerOpen = useUserModalStoreIsDetailedUserInfoDrawerOpen()
    const selectedUser = useUserModalStoreDrawerUserUuid()

    const cleanUpDrawer = async () => {
        actions.changeDetailedUserInfoDrawerState(false)
    }

    const isQueryEnabled = !!selectedUser

    const { data: user, isLoading: isUserLoading } = useGetUserByUuid({
        route: {
            uuid: selectedUser ?? ''
        },
        rQueryParams: {
            enabled: isQueryEnabled
        }
    })
    const { data: subNodes } = useGetSubscriptionConnections()

    const subscriptionPageLinks = user
        ? buildSubscriptionLinksFromNodes(
              (subNodes ?? []).filter((node) => !node.isDisabled),
              user.shortUuid,
              user.subscriptionUrl
          )
        : []

    return (
        <Drawer
            keepMounted={true}
            onClose={cleanUpDrawer}
            opened={isDetailedUserInfoDrawerOpen}
            overlayProps={{ backgroundOpacity: 0.6, blur: 0 }}
            padding="lg"
            position="right"
            size="lg"
            title={
                <BaseOverlayHeader
                    iconColor="blue"
                    IconComponent={PiUserDuotone}
                    iconVariant="soft"
                    title={t('detailed-user-info-drawer.widget.detailed-user-info')}
                />
            }
            zIndex={1000}
        >
            {isUserLoading && (
                <Center h="100%" mt="md" py="xl" ta="center">
                    <Box>
                        <LoaderModalShared
                            text={t('detailed-user-info-drawer.widget.loading-user-info')}
                        />
                    </Box>
                </Center>
            )}

            {!isUserLoading && user && (
                <Stack gap="md">
                    <SectionCardRoot>
                        <SectionCardSection>
                            <Group align="center" justify="space-between">
                                <BaseOverlayHeader
                                    iconColor="blue"
                                    IconComponent={PiUserDuotone}
                                    iconVariant="soft"
                                    title={t('detailed-user-info-drawer.widget.user-information')}
                                />

                                <Group gap="xs">
                                    <UserStatusBadge
                                        h={28}
                                        key="view-user-status-badge"
                                        size="lg"
                                        status={user.status}
                                    />
                                </Group>
                            </Group>
                        </SectionCardSection>
                        <SectionCardSection>
                            <DataList withDivider orientation="vertical">
                                <CopyableDataListItem label="ID" monospace value={user.id} />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.uuid')}
                                    monospace
                                    value={user.uuid}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.short-uuid')}
                                    monospace
                                    value={user.shortUuid}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.username')}
                                    value={user.username}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.email')}
                                    value={user.email}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.telegram-id')}
                                    monospace
                                    value={user.telegramId}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.description')}
                                    value={user.description}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.tag')}
                                    value={user.tag}
                                />
                            </DataList>
                        </SectionCardSection>
                    </SectionCardRoot>
                    <SectionCardRoot>
                        <SectionCardSection>
                            <BaseOverlayHeader
                                iconColor="teal"
                                IconComponent={PiArrowsDownUpDuotone}
                                iconVariant="soft"
                                title={t('detailed-user-info-drawer.widget.traffic-information')}
                            />
                        </SectionCardSection>
                        <SectionCardSection>
                            <DataList withDivider orientation="vertical">
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.used-traffic')}
                                    value={prettifyBytesUtil(user.userTraffic.usedTrafficBytes)}
                                />
                                <CopyableDataListItem
                                    label={t(
                                        'detailed-user-info-drawer.widget.lifetime-used-traffic'
                                    )}
                                    value={prettifyBytesUtil(
                                        user.userTraffic.lifetimeUsedTrafficBytes
                                    )}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.traffic-limit')}
                                    value={prettifyBytesUtil(user.trafficLimitBytes)}
                                />
                                <CopyableDataListItem
                                    label={t(
                                        'detailed-user-info-drawer.widget.traffic-limit-strategy'
                                    )}
                                    value={user.trafficLimitStrategy}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.last-traffic-reset')}
                                    value={formatTimeUtil({
                                        time: user.lastTrafficResetAt,
                                        template: 'TIME_FIRST_DATETIME',
                                        language: i18n.language
                                    })}
                                />
                            </DataList>
                        </SectionCardSection>
                    </SectionCardRoot>
                    <SectionCardRoot>
                        <SectionCardSection>
                            <BaseOverlayHeader
                                iconColor="orange"
                                IconComponent={PiCalendarDotDuotone}
                                iconSize={16}
                                iconVariant="soft"
                                title={t(
                                    'detailed-user-info-drawer.widget.subscription-information'
                                )}
                            />
                        </SectionCardSection>
                        <SectionCardSection>
                            <DataList withDivider orientation="vertical">
                                {subscriptionPageLinks.map((link, index) => (
                                    <CopyableDataListItem
                                        key={link.url}
                                        label={
                                            index === 0
                                                ? t(
                                                      'detailed-user-info-drawer.widget.subscription-url'
                                                  )
                                                : link.nodeName
                                        }
                                        monospace
                                        value={link.url}
                                    />
                                ))}
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.expires-at')}
                                    value={formatTimeUtil({
                                        time: user.expireAt,
                                        template: 'TIME_FIRST_DATETIME',
                                        language: i18n.language
                                    })}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.revoked-at')}
                                    value={formatTimeUtil({
                                        time: user.subRevokedAt,
                                        template: 'TIME_FIRST_DATETIME',
                                        language: i18n.language
                                    })}
                                />
                            </DataList>
                        </SectionCardSection>
                    </SectionCardRoot>

                    <SectionCardRoot>
                        <SectionCardSection>
                            <BaseOverlayHeader
                                iconColor="violet"
                                IconComponent={PiNetworkDuotone}
                                iconVariant="soft"
                                title={t('detailed-user-info-drawer.widget.connection-information')}
                            />
                        </SectionCardSection>
                        <SectionCardSection>
                            <DataList withDivider orientation="vertical">
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.trojan-password')}
                                    monospace
                                    value={user.trojanPassword}
                                />
                                <CopyableDataListItem
                                    label="VLESS UUID"
                                    monospace
                                    value={user.vlessUuid}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.ss-password')}
                                    monospace
                                    value={user.ssPassword}
                                />
                                {(() => {
                                    const protocolCredentials = user as typeof user & {
                                        naivePassword?: string
                                        shadowtlsPassword?: string
                                        hysteria2Password?: string
                                        anytlsPassword?: string
                                    }

                                    return (
                                        <>
                                            <CopyableDataListItem
                                                label="Naive Password"
                                                monospace
                                                value={protocolCredentials.naivePassword}
                                            />
                                            <CopyableDataListItem
                                                label="ShadowTLS Password"
                                                monospace
                                                value={protocolCredentials.shadowtlsPassword}
                                            />
                                            <CopyableDataListItem
                                                label="Hysteria2 Password"
                                                monospace
                                                value={protocolCredentials.hysteria2Password}
                                            />
                                            <CopyableDataListItem
                                                label="AnyTLS Password"
                                                monospace
                                                value={protocolCredentials.anytlsPassword}
                                            />
                                        </>
                                    )
                                })()}
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.first-connected-at')}
                                    value={formatTimeUtil({
                                        time: user.userTraffic.firstConnectedAt,
                                        template: 'TIME_FIRST_DATETIME',
                                        language: i18n.language
                                    })}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.last-online')}
                                    value={formatTimeUtil({
                                        time: user.userTraffic.onlineAt,
                                        template: 'TIME_FIRST_DATETIME',
                                        language: i18n.language
                                    })}
                                />
                                <CopyableDataListItem
                                    label={t(
                                        'detailed-user-info-drawer.widget.last-connected-node'
                                    )}
                                    monospace
                                    value={user.userTraffic.lastConnectedNodeUuid}
                                />
                            </DataList>
                        </SectionCardSection>
                    </SectionCardRoot>

                    {user.activeInternalSquads && user.activeInternalSquads.length > 0 && (
                        <SectionCardRoot>
                            <SectionCardSection>
                                <BaseOverlayHeader
                                    iconColor="green"
                                    IconComponent={PiTagDuotone}
                                    iconVariant="soft"
                                    title={t(
                                        'detailed-user-info-drawer.widget.active-internal-squads'
                                    )}
                                />
                            </SectionCardSection>
                            <SectionCardSection>
                                <DataList withDivider orientation="vertical">
                                    {user.activeInternalSquads.map((squad) => (
                                        <CopyableDataListItem
                                            key={squad.uuid}
                                            label={squad.name}
                                            monospace
                                            value={squad.uuid}
                                        />
                                    ))}
                                </DataList>
                            </SectionCardSection>
                        </SectionCardRoot>
                    )}
                    <SectionCardRoot>
                        <SectionCardSection>
                            <BaseOverlayHeader
                                iconColor="gray"
                                IconComponent={PiClockDuotone}
                                iconVariant="soft"
                                title={t('detailed-user-info-drawer.widget.timestamps')}
                            />
                        </SectionCardSection>
                        <SectionCardSection>
                            <DataList withDivider orientation="vertical">
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.created-at')}
                                    value={formatTimeUtil({
                                        time: user.createdAt,
                                        template: 'TIME_FIRST_DATETIME',
                                        language: i18n.language
                                    })}
                                />
                                <CopyableDataListItem
                                    label={t('detailed-user-info-drawer.widget.updated-at')}
                                    value={formatTimeUtil({
                                        time: user.updatedAt,
                                        template: 'TIME_FIRST_DATETIME',
                                        language: i18n.language
                                    })}
                                />
                            </DataList>
                        </SectionCardSection>
                    </SectionCardRoot>
                </Stack>
            )}
        </Drawer>
    )
}
