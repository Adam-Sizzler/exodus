import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import {
    ActionIcon,
    Badge,
    Box,
    Checkbox,
    Group,
    OverflowList,
    px,
    Stack,
    Text,
    ThemeIcon,
    Tooltip
} from '@mantine/core'
import {
    GetAllHostsCommand,
    GetAllNodesCommand,
    GetConfigProfilesCommand
} from '@exodus/backend-contract'
import cx from 'clsx'
import ColorHash from 'color-hash'
import { CSSProperties, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { PiNetwork, PiProhibit, PiPulse } from 'react-icons/pi'
import { RiDraggable } from 'react-icons/ri'
import {
    TbAlertCircle,
    TbCirclesRelation,
    TbCloudNetwork,
    TbEyeOff,
    TbFileDescription,
    TbMask,
    TbStar
} from 'react-icons/tb'
import { createSearchParams, useNavigate } from 'react-router'

import { ROUTES } from '@shared/constants'
import { SEARCH_PARAMS } from '@shared/constants/search-params'
import { useIsMobile } from '@shared/hooks'
import { SingboxLogo } from '@shared/ui/logos'
import { resolveCountryCode } from '@shared/utils/misc/resolve-country-code'
import { openOrNavigate } from '@shared/utils/open-or-navigate'

import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'

import classes from './HostCard.module.css'

export interface IProps {
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles'] | undefined
    isDragOverlay?: boolean
    isSelected?: boolean
    item: GetAllHostsCommand.Response['response'][number]
    nodesByUuid: Map<string, GetAllNodesCommand.Response['response'][number]>
    onSelect?: () => void
    openExternal?: boolean
    viewOnly?: boolean
    disableReordering?: boolean
}

export function HostCardWidget(props: IProps) {
    const {
        nodesByUuid,
        item,
        configProfiles,
        isSelected,
        onSelect,
        isDragOverlay = false,
        viewOnly = false,
        openExternal = false,
        disableReordering = false
    } = props

    const { t } = useTranslation()
    const navigate = useNavigate()

    const openModalWithData = useModalsStoreOpenWithData()

    const [isHovered, setIsHovered] = useState(false)
    const isMobile = useIsMobile()

    const configProfile = configProfiles?.find(
        (configProfile) => configProfile.uuid === item.inbound.configProfileUuid
    )

    const inboundTag = configProfile?.inbounds.find(
        (inbound) => inbound.uuid === item.inbound.configProfileInboundUuid
    )?.tag

    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
        id: item.uuid,
        disabled: isDragOverlay
    })

    const style: CSSProperties = {
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0 : 1,
        zIndex: isDragging ? 1000 : 'auto',
        position: 'relative'
    }

    const handleEdit = () => {
        if (openExternal) {
            openOrNavigate(
                `${ROUTES.DASHBOARD.MANAGEMENT.HOSTS}?${createSearchParams({
                    [SEARCH_PARAMS.HOST]: item.uuid
                })}`,
                navigate
            )

            return
        }

        openModalWithData(MODALS.EDIT_HOST_MODAL, item)
    }

    if (!configProfiles) {
        return null
    }

    const isHostActive = !item.isDisabled

    const ch = new ColorHash({ lightness: [0.65, 0.65, 0.65] })

    const isParamSet = (value: unknown): boolean => {
        if (value == null) return false
        if (typeof value === 'string') return value.trim().length > 0
        if (Array.isArray(value)) return value.length > 0
        if (typeof value === 'object') return Object.keys(value as object).length > 0
        return true
    }

    const hasFinalMask = isParamSet(item.finalMask)
    const hasMuxParams = isParamSet(item.muxParams)
    const hasSockoptParams = isParamSet(item.sockoptParams)
    const hasXrayJsonTemplate = !!item.xrayJsonTemplateUuid
    const serverDescription = item.serverDescription?.trim() || ''
    const hasExcludedSquads = item.excludedInternalSquads.length > 0

    if (isMobile) {
        return (
            <Box
                className={cx(classes.item, classes.mobileItem, {
                    [classes.itemDragging]: isDragging || isHovered,
                    [classes.selectedItem]: isSelected,
                    [classes.danglingItem]: !configProfile?.uuid
                })}
                data-dnd-overlay={isDragOverlay}
                ref={isDragOverlay ? undefined : setNodeRef}
                style={style}
            >
                <Stack gap="sm">
                    <Group justify="space-between" wrap="nowrap">
                        {!viewOnly && (
                            <Group gap="sm" wrap="nowrap">
                                <Checkbox
                                    checked={isSelected}
                                    onChange={(e) => {
                                        e.stopPropagation()
                                        onSelect?.()
                                    }}
                                    size="md"
                                    styles={{
                                        input: { cursor: 'pointer' }
                                    }}
                                />
                                {!disableReordering && (
                                    <Box
                                        {...(isDragOverlay ? {} : attributes)}
                                        {...(isDragOverlay ? {} : listeners)}
                                        className={classes.mobileDragHandle}
                                        onClick={(e) => e.stopPropagation()}
                                    >
                                        <RiDraggable size={px('1.2rem')} />
                                    </Box>
                                )}
                            </Group>
                        )}

                        {!isHostActive && (
                            <ActionIcon
                                color="gray"
                                size="lg"
                                style={{ flexShrink: 0 }}
                                variant="soft"
                            >
                                <PiProhibit size={16} />
                            </ActionIcon>
                        )}

                        {isHostActive && item.isHidden && (
                            <ActionIcon
                                color="violet"
                                size="lg"
                                style={{ flexShrink: 0 }}
                                variant="soft"
                            >
                                <TbEyeOff size={16} />
                            </ActionIcon>
                        )}

                        {isHostActive && !item.isHidden && (
                            <ActionIcon
                                color="teal"
                                size="lg"
                                style={{ flexShrink: 0 }}
                                variant="soft"
                            >
                                <PiPulse size={16} />
                            </ActionIcon>
                        )}
                    </Group>

                    <Box
                        className={classes.contentArea}
                        onClick={(e) => {
                            e.stopPropagation()
                            handleEdit()
                        }}
                        onTouchEnd={() => setIsHovered(false)}
                        onTouchStart={() => setIsHovered(true)}
                    >
                        <Stack gap="xs">
                            <Text className={classes.hostName} fw={600} size="lg">
                                {item.remark}
                            </Text>

                            <Text c="dimmed" className={classes.mobileAddress} size="sm">
                                {item.address}
                                {item.port ? `:${item.port}` : ''}
                            </Text>

                            <Stack gap="xs">
                                <Badge
                                    autoContrast
                                    color={configProfile?.uuid ? ch.hex(configProfile.uuid) : 'red'}
                                    leftSection={
                                        configProfile?.uuid ? (
                                            <SingboxLogo size={12} />
                                        ) : (
                                            <TbAlertCircle size={12} />
                                        )
                                    }
                                    size="md"
                                    style={{
                                        paddingLeft: 0
                                    }}
                                    variant="transparent"
                                >
                                    {configProfile?.name || 'DANGLING'}
                                    {item.inbound.configProfileInboundUuid && (
                                        <>
                                            <span style={{ margin: '0 6px', opacity: 0.5 }}>›</span>
                                            <span style={{ opacity: 0.75 }}>
                                                {inboundTag || 'UNKNOWN'}
                                            </span>
                                        </>
                                    )}
                                </Badge>

                                <OverflowList
                                    data={item.tags.sort((a, b) => a.localeCompare(b))}
                                    gap={0}
                                    maxRows={1}
                                    maxVisibleItems={2}
                                    renderItem={(tag) => (
                                        <Badge
                                            autoContrast
                                            color={ch.hex(tag)}
                                            key={tag}
                                            leftSection={<TbStar size={12} />}
                                            size="md"
                                            style={{
                                                paddingLeft: 0
                                            }}
                                            variant="transparent"
                                        >
                                            {tag}
                                        </Badge>
                                    )}
                                    renderOverflow={(items) => (
                                        <Tooltip
                                            label={
                                                <Stack gap="xs">
                                                    {items.map((tag) => (
                                                        <Badge
                                                            color={ch.hex(tag)}
                                                            fullWidth
                                                            key={tag}
                                                            leftSection={<TbStar size={12} />}
                                                            size="md"
                                                            variant="transparent"
                                                        >
                                                            {tag}
                                                        </Badge>
                                                    ))}
                                                </Stack>
                                            }
                                            multiline
                                            position="top"
                                        >
                                            <Badge color="gray" size="md" variant="transparent">
                                                +{items.length}
                                            </Badge>
                                        </Tooltip>
                                    )}
                                />
                            </Stack>
                        </Stack>
                    </Box>
                </Stack>
            </Box>
        )
    }

    return (
        <Box
            className={cx(classes.item, {
                [classes.itemDragging]: isDragging || isHovered,
                [classes.selectedItem]: isSelected,
                [classes.danglingItem]: !configProfile?.uuid
            })}
            data-dnd-overlay={isDragOverlay}
            ref={isDragOverlay ? undefined : setNodeRef}
            style={style}
        >
            <Group gap="md" w="100%" wrap="nowrap">
                {!viewOnly && (
                    <Group gap="xs" wrap="nowrap">
                        <Checkbox checked={isSelected} onChange={onSelect} size="md" />
                        {!disableReordering && (
                            <Box
                                {...(isDragOverlay ? {} : attributes)}
                                {...(isDragOverlay ? {} : listeners)}
                                className={classes.dragHandle}
                            >
                                <RiDraggable color="white" size="24px" />
                            </Box>
                        )}
                    </Group>
                )}

                <Stack
                    className={classes.contentArea}
                    flex={1}
                    gap={6}
                    miw={0}
                    onClick={handleEdit}
                    onMouseEnter={() => setIsHovered(true)}
                    onMouseLeave={() => setIsHovered(false)}
                >
                    <Group
                        gap="sm"
                        justify="space-between"
                        miw={0}
                        style={{ overflow: 'hidden' }}
                        wrap="nowrap"
                    >
                        <Group
                            gap="sm"
                            miw={0}
                            style={{ flexShrink: 1, overflow: 'hidden' }}
                            wrap="nowrap"
                        >
                            {!isHostActive && (
                                <ActionIcon
                                    color="gray"
                                    size="md"
                                    style={{ flexShrink: 0 }}
                                    variant="soft"
                                >
                                    <PiProhibit size={16} />
                                </ActionIcon>
                            )}

                            {isHostActive && item.isHidden && (
                                <ActionIcon
                                    color="violet"
                                    size="md"
                                    style={{ flexShrink: 0 }}
                                    variant="soft"
                                >
                                    <TbEyeOff size={16} />
                                </ActionIcon>
                            )}

                            {isHostActive && !item.isHidden && (
                                <ActionIcon
                                    color="teal"
                                    size="md"
                                    style={{ flexShrink: 0 }}
                                    variant="soft"
                                >
                                    <PiPulse size={16} />
                                </ActionIcon>
                            )}

                            <Text fw={600} style={{ flexShrink: 0 }} truncate>
                                {item.remark}
                            </Text>
                            <Text
                                c="dimmed"
                                className={classes.hostAddress}
                                style={{ flexShrink: 1, minWidth: 0 }}
                                truncate
                            >
                                {item.address}
                                {item.port ? `:${item.port}` : ''}
                            </Text>
                        </Group>

                        <Group gap={6} style={{ flexShrink: 0 }} wrap="nowrap">
                            {hasExcludedSquads && (
                                <Tooltip label={t('base-host-form.excluded-internal-squads')}>
                                    <ThemeIcon color="yellow" size={28} variant="soft">
                                        <TbCirclesRelation size={16} />
                                    </ThemeIcon>
                                </Tooltip>
                            )}

                            <Tooltip label={t('base-host-form.xray-json-template')}>
                                <ThemeIcon
                                    color={hasXrayJsonTemplate ? 'teal' : 'gray'}
                                    size={28}
                                    variant="soft"
                                >
                                    <SingboxLogo size={16} />
                                </ThemeIcon>
                            </Tooltip>

                            <Tooltip label="Mux">
                                <ThemeIcon
                                    color={hasMuxParams ? 'teal' : 'gray'}
                                    size={28}
                                    variant="soft"
                                >
                                    <TbCloudNetwork size={16} />
                                </ThemeIcon>
                            </Tooltip>

                            <Tooltip label="Final Mask">
                                <ThemeIcon
                                    color={hasFinalMask ? 'teal' : 'gray'}
                                    size={28}
                                    variant="soft"
                                >
                                    <TbMask size={16} />
                                </ThemeIcon>
                            </Tooltip>

                            <Tooltip label="SockOpt">
                                <ThemeIcon
                                    color={hasSockoptParams ? 'teal' : 'gray'}
                                    size={28}
                                    variant="soft"
                                >
                                    <PiNetwork size={16} />
                                </ThemeIcon>
                            </Tooltip>
                        </Group>
                    </Group>

                    <Group
                        gap="xs"
                        justify="space-between"
                        miw={0}
                        style={{ overflow: 'hidden' }}
                        wrap="nowrap"
                    >
                        <Group
                            gap={0}
                            miw={0}
                            style={{ flexShrink: 1, overflow: 'hidden' }}
                            wrap="nowrap"
                        >
                            <Badge
                                autoContrast
                                color={configProfile?.uuid ? ch.hex(configProfile.uuid) : 'red'}
                                leftSection={
                                    configProfile?.uuid ? (
                                        <SingboxLogo size={12} />
                                    ) : (
                                        <TbAlertCircle size={12} />
                                    )
                                }
                                size="md"
                                style={{
                                    paddingLeft: 0
                                }}
                                variant="transparent"
                            >
                                {configProfile?.name || 'DANGLING'}
                                {item.inbound.configProfileInboundUuid && (
                                    <>
                                        <span style={{ margin: '0 6px', opacity: 0.5 }}>›</span>
                                        <span style={{ opacity: 0.75 }}>
                                            {inboundTag || 'UNKNOWN'}
                                        </span>
                                    </>
                                )}
                            </Badge>

                            <OverflowList
                                data={item.tags.sort((a, b) => a.localeCompare(b))}
                                gap={0}
                                maxRows={1}
                                maxVisibleItems={2}
                                renderItem={(tag) => (
                                    <Badge
                                        autoContrast
                                        color={ch.hex(tag)}
                                        key={tag}
                                        leftSection={<TbStar size={12} />}
                                        size="md"
                                        variant="transparent"
                                    >
                                        {tag}
                                    </Badge>
                                )}
                                renderOverflow={(items) => (
                                    <Tooltip
                                        label={
                                            <Stack gap="xs">
                                                {items.map((tag) => (
                                                    <Badge
                                                        color={ch.hex(tag)}
                                                        fullWidth
                                                        key={tag}
                                                        leftSection={<TbStar size={12} />}
                                                        size="md"
                                                        variant="transparent"
                                                    >
                                                        {tag}
                                                    </Badge>
                                                ))}
                                            </Stack>
                                        }
                                        multiline
                                        position="top"
                                    >
                                        <Badge color="gray" size="md" variant="transparent">
                                            +{items.length}
                                        </Badge>
                                    </Tooltip>
                                )}
                            />

                            {serverDescription && (
                                <Badge
                                    color="dark.3"
                                    leftSection={<TbFileDescription size={12} />}
                                    size="md"
                                    style={{
                                        fontStyle: 'italic',
                                        fontWeight: 500,
                                        letterSpacing: '1px'
                                    }}
                                    variant="transparent"
                                >
                                    {serverDescription}
                                </Badge>
                            )}
                        </Group>

                        <Group gap="xs" style={{ flexShrink: 0 }} wrap="nowrap">
                            <OverflowList
                                data={item.nodes
                                    .map((nodeId) => nodesByUuid.get(nodeId))
                                    .filter((n): n is NonNullable<typeof n> => Boolean(n))}
                                gap={4}
                                maxRows={1}
                                maxVisibleItems={3}
                                renderItem={(node) => (
                                    <Badge
                                        autoContrast
                                        color="gray"
                                        key={`${node.uuid}|${item.uuid}`}
                                        leftSection={resolveCountryCode(node.countryCode, 18)}
                                        size="md"
                                        style={{ cursor: 'pointer' }}
                                        variant="soft"
                                    >
                                        {node.name}
                                    </Badge>
                                )}
                                renderOverflow={(items) => (
                                    <Tooltip
                                        label={
                                            <Stack gap="xs">
                                                {items.map((node) => (
                                                    <Badge
                                                        color="gray"
                                                        fullWidth
                                                        key={`${node.uuid}|${item.uuid}`}
                                                        leftSection={resolveCountryCode(
                                                            node.countryCode,
                                                            18
                                                        )}
                                                        variant="soft"
                                                    >
                                                        {node.name}
                                                    </Badge>
                                                ))}
                                            </Stack>
                                        }
                                        multiline
                                        position="top"
                                    >
                                        <Badge color="gray" size="md" variant="soft">
                                            +{items.length}
                                        </Badge>
                                    </Tooltip>
                                )}
                            />
                        </Group>
                    </Group>
                </Stack>
            </Group>
        </Box>
    )
}
