import { PiDotsSixVertical, PiGlobeSimple } from 'react-icons/pi'
import { Avatar, Badge, Box, Flex, Grid, Text } from '@mantine/core'
import { notifications } from '@mantine/notifications'
import ReactCountryFlag from 'react-country-flag'
import { useSortable } from '@dnd-kit/sortable'
import { useTranslation } from 'react-i18next'
import { useClipboard } from '@mantine/hooks'
import { CSSProperties, memo } from 'react'
import { CSS } from '@dnd-kit/utilities'
import clsx from 'clsx'

import { getSingboxUptimeUtil } from '@shared/utils/time-utils'
import { faviconResolver } from '@shared/utils/misc'
import { Logo } from '@shared/ui'

import { NodeStatusBadgeWidget } from '../node-status-badge'
import classes from './NodeCard.module.css'
import { IProps } from './interfaces'

const getNodeColors = (node: IProps['node']) => {
    if (node.isDisabled) {
        return {
            backgroundColor: 'rgba(107, 114, 128, 0.15)',
            borderColor: 'rgba(107, 114, 128, 0.3)',
            boxShadow: 'rgba(107, 114, 128, 0.2)'
        }
    }
    if (node.isConnected) {
        return {
            backgroundColor: 'rgba(45, 212, 191, 0.15)',
            borderColor: 'rgba(45, 212, 191, 0.3)',
            boxShadow: 'rgba(45, 212, 191, 0.2)'
        }
    }
    if (node.isConnecting) {
        return {
            backgroundColor: 'rgba(245, 158, 11, 0.15)',
            borderColor: 'rgba(245, 158, 11, 0.3)',
            boxShadow: 'rgba(245, 158, 11, 0.2)'
        }
    }
    return {
        backgroundColor: 'rgba(239, 68, 68, 0.15)',
        borderColor: 'rgba(239, 68, 68, 0.3)',
        boxShadow: 'rgba(239, 68, 68, 0.2)'
    }
}

export const NodeCardWidget = memo((props: IProps) => {
    const { t } = useTranslation()
    const { handleViewNode, node, isDragOverlay = false, isMobile } = props

    const clipboard = useClipboard({ timeout: 500 })

    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
        id: node.uuid
    })

    const style: CSSProperties = {
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0 : 1,
        zIndex: isDragging ? 1000 : 'auto'
    }

    const nodeSingboxUptime = node.singboxUptime
    const isOnline = node.isConnected && nodeSingboxUptime !== 0 && !node.isDisabled
    const { backgroundColor, borderColor, boxShadow } = getNodeColors(node)

    const handleCopy = (e: React.MouseEvent) => {
        e.stopPropagation()
        clipboard.copy(node.address)
        notifications.show({
            message: node.address,
            title: t('node-card.widget.copied'),
            color: 'teal'
        })
    }

    return (
        <Box
            className={clsx(classes.nodeRow, {
                [classes.nodeRowDragging]: isDragging
            })}
            data-dnd-overlay={isDragOverlay}
            onClick={() => handleViewNode(node.uuid)}
            ref={isDragOverlay ? undefined : setNodeRef}
            style={{
                ...style,
                background: `linear-gradient(
                    135deg,
                    ${backgroundColor} 0%,
                    var(--mantine-color-dark-7) 100%
                )`,
                borderColor,
                boxShadow
            }}
        >
            <Box
                {...(isDragOverlay ? {} : attributes)}
                {...(isDragOverlay ? {} : listeners)}
                className={clsx(classes.dragHandle, {
                    [classes.dragHandleActive]: isDragging
                })}
            >
                <PiDotsSixVertical color="white" size="24px" />
            </Box>

            {!isMobile && (
                <Grid align="center" className={classes.desktopGrid} gutter="md">
                    <Grid.Col span={{ base: 12, sm: 6.5 }}>
                        <Flex align="center" gap="sm">
                            <NodeStatusBadgeWidget node={node} withText={false} />

                            <Flex align="center" className={classes.nameContainer} gap="xs">
                                {node.countryCode && node.countryCode !== 'XX' && (
                                    <ReactCountryFlag
                                        countryCode={node.countryCode}
                                        style={{
                                            fontSize: '1.6em',
                                            borderRadius: '2px'
                                        }}
                                    />
                                )}
                                <Text className={classes.nodeName} fw={600} size="md">
                                    {node.name}
                                </Text>
                            </Flex>

                            <Flex align="center" gap="xs">
                                {node.provider && (
                                    <Badge
                                        color="gray"
                                        leftSection={
                                            <Avatar
                                                alt={node.provider.name}
                                                color="initials"
                                                name={node.provider.name}
                                                onLoad={(event) => {
                                                    const img = event.target as HTMLImageElement
                                                    if (
                                                        img.naturalWidth <= 16 &&
                                                        img.naturalHeight <= 16
                                                    ) {
                                                        img.src = ''
                                                    }
                                                }}
                                                radius="sm"
                                                size={16}
                                                src={faviconResolver(node.provider.faviconLink)}
                                            />
                                        }
                                        size="lg"
                                        style={{
                                            maxWidth: '20ch',
                                            overflow: 'hidden',
                                            textOverflow: 'ellipsis',
                                            whiteSpace: 'nowrap',
                                            cursor: 'pointer'
                                        }}
                                        variant="light"
                                    >
                                        {node.provider.name}
                                    </Badge>
                                )}
                            </Flex>
                        </Flex>
                    </Grid.Col>

                    <Grid.Col span={{ base: 12, sm: 3 }}>
                        <Flex align="center" gap="xs">
                            <PiGlobeSimple className={classes.icon} size={14} />
                            <Text
                                c="dimmed"
                                className={classes.addressText}
                                onClick={handleCopy}
                                size="sm"
                            >
                                {node.address}
                            </Text>
                        </Flex>
                    </Grid.Col>


                    <Grid.Col span={{ base: 12, sm: 2.5 }}>
                        <Flex align="center" gap="xs" justify="flex-end">
                            {isOnline && (
                                <Flex align="center" gap={4}>
                                    <Logo size={14} />
                                    <Text c="teal" fw={600} size="sm" truncate>
                                        {getSingboxUptimeUtil(nodeSingboxUptime)}
                                    </Text>
                                </Flex>
                            )}
                        </Flex>
                    </Grid.Col>
                </Grid>
            )}

            {isMobile && (
                <Box>
                    <Flex align="center" gap="sm" mb="xs">
                        <NodeStatusBadgeWidget node={node} withText={false} />

                        <Flex align="center" gap="xs" style={{ flex: 1, minWidth: 0 }}>
                            {node.countryCode && node.countryCode !== 'XX' && (
                                <ReactCountryFlag
                                    countryCode={node.countryCode}
                                    style={{
                                        fontSize: '1.5em',
                                        borderRadius: '2px'
                                    }}
                                />
                            )}
                            <Text className={classes.nodeName} fw={600} size="sm">
                                {node.name}
                            </Text>
                        </Flex>
                    </Flex>

                    <Box mb="xs">
                        <Flex align="center" gap="xs">
                            {node.provider && (
                                <Badge
                                    color="gray"
                                    leftSection={
                                        <Avatar
                                            alt={node.provider.name}
                                            color="initials"
                                            name={node.provider.name}
                                            onLoad={(event) => {
                                                const img = event.target as HTMLImageElement
                                                if (
                                                    img.naturalWidth <= 16 &&
                                                    img.naturalHeight <= 16
                                                ) {
                                                    img.src = ''
                                                }
                                            }}
                                            radius="sm"
                                            size={16}
                                            src={faviconResolver(node.provider.faviconLink)}
                                        />
                                    }
                                    size="lg"
                                    variant="light"
                                >
                                    {node.provider.name}
                                </Badge>
                            )}

                            {!node.provider && (
                                <Badge
                                    color="gray"
                                    leftSection={
                                        <Avatar
                                            alt="Unknown"
                                            color="initials"
                                            name="Unknown"
                                            radius="sm"
                                            size={16}
                                        />
                                    }
                                    size="lg"
                                    style={{
                                        visibility: 'hidden'
                                    }}
                                    variant="light"
                                >
                                    Unknown
                                </Badge>
                            )}
                        </Flex>
                    </Box>


                    <Flex align="center" justify="flex-end" mt="xs">
                        <Flex align="center" gap={4}>
                            <Logo size={12} />
                            <Text
                                c={isOnline ? 'teal' : 'dimmed'}
                                fw={isOnline ? 600 : 500}
                                size="xs"
                            >
                                {isOnline ? getSingboxUptimeUtil(nodeSingboxUptime) : 'offline'}
                            </Text>
                        </Flex>
                    </Flex>
                </Box>
            )}
        </Box>
    )
})
