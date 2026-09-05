import { Badge, Checkbox, CopyButton, Group, Menu, Tooltip } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import {
    PiCheck,
    PiClock,
    PiCopy,
    PiListChecks,
    PiPencil,
    PiProhibit,
    PiPulse,
    PiTrashDuotone,
    PiWarningCircle
} from 'react-icons/pi'
import { TbLink, TbTags } from 'react-icons/tb'

import { showModal } from '@shared/_modals/show-modal'
import { WithDndSortable } from '@shared/hocs/with-dnd-sortable'
import { EntityCardShared } from '@shared/ui/entity-card'

export type SRSListItem = {
    uuid: string
    url: string
    updateInterval: string
    fileName: string
    shortName: string
    viewPosition: number
    isEnabled: boolean
    isAvailable: boolean
    lastError?: null | string
    tags?: string[]
}

interface IProps {
    disableReordering?: boolean
    isDragOverlay?: boolean
    isSelected?: boolean
    onCheck: (uuid: string) => void
    onDelete: (uuid: string, name: string) => void
    onEdit: (item: SRSListItem) => void
    onToggleEnable: (uuid: string, enabled: boolean) => void
    onToggleSelect?: (uuid: string) => void
    srsList: SRSListItem
}

export function SRSListCardWidget(props: IProps) {
    const {
        disableReordering = false,
        isDragOverlay = false,
        isSelected = false,
        onCheck,
        onDelete,
        onEdit,
        onToggleEnable,
        onToggleSelect,
        srsList
    } = props

    const { t } = useTranslation()
    const tr = (key: string, defaultValue: string) => t(key, { defaultValue })

    const isActive = srsList.isEnabled && srsList.isAvailable
    const displayName = srsList.fileName || srsList.shortName || srsList.url
    const stateColor = !srsList.isEnabled ? 'gray' : srsList.isAvailable ? 'teal' : 'red'

    return (
        <WithDndSortable
            disableReordering={disableReordering}
            dragHandlePosition="inline-end"
            id={srsList.uuid}
            isDragOverlay={isDragOverlay}
        >
            <EntityCardShared.Root isActive={isActive} onClick={() => onEdit(srsList)}>
                <EntityCardShared.Header>
                    {onToggleSelect && (
                        <Checkbox
                            checked={isSelected}
                            onClick={(event) => event.stopPropagation()}
                            onChange={() => onToggleSelect(srsList.uuid)}
                            size="sm"
                            style={{ cursor: 'pointer', flexShrink: 0 }}
                        />
                    )}
                    <EntityCardShared.Icon highlight={isActive} color={stateColor}>
                        {!srsList.isEnabled ? (
                            <PiProhibit size={20} />
                        ) : srsList.isAvailable ? (
                            <PiListChecks size={22} />
                        ) : (
                            <PiWarningCircle size={22} />
                        )}
                    </EntityCardShared.Icon>
                    <EntityCardShared.Content
                        tags={srsList.tags ?? []}
                        badges={
                            <Group gap="xs" wrap="nowrap">
                                <Tooltip
                                    label={tr(
                                        'srs-lists.feature.update-interval',
                                        'Update interval'
                                    )}
                                    withArrow
                                >
                                    <Badge
                                        color="blue"
                                        leftSection={<PiClock size={12} />}
                                        size="lg"
                                        variant="soft"
                                    >
                                        {srsList.updateInterval}
                                    </Badge>
                                </Tooltip>

                                <Tooltip
                                    label={
                                        !srsList.isEnabled
                                            ? tr('srs-lists.feature.disabled', 'Disabled')
                                            : srsList.isAvailable
                                              ? tr('srs-lists.feature.available', 'Available')
                                              : srsList.lastError
                                                ? `${tr('srs-lists.feature.unavailable', 'Unavailable')}. ${tr('srs-lists.feature.last-error', 'Error')}: ${srsList.lastError}`
                                                : tr('srs-lists.feature.unavailable', 'Unavailable')
                                    }
                                    maw={420}
                                    multiline
                                    withArrow
                                >
                                    <Badge
                                        color={stateColor}
                                        leftSection={
                                            !srsList.isEnabled ? (
                                                <PiProhibit size={12} />
                                            ) : srsList.isAvailable ? (
                                                <PiPulse size={12} />
                                            ) : (
                                                <PiWarningCircle size={12} />
                                            )
                                        }
                                        size="lg"
                                        variant="soft"
                                    >
                                        {!srsList.isEnabled
                                            ? tr('srs-lists.feature.disabled', 'Disabled')
                                            : srsList.isAvailable
                                              ? tr('srs-lists.feature.available', 'Available')
                                              : tr('srs-lists.feature.unavailable', 'Unavailable')}
                                    </Badge>
                                </Tooltip>
                            </Group>
                        }
                        title={displayName}
                    />
                </EntityCardShared.Header>

                <EntityCardShared.Actions>
                    <EntityCardShared.Menu>
                        <Menu.Item
                            leftSection={<PiPencil size={18} />}
                            onClick={() => onEdit(srsList)}
                        >
                            {tr('common.action.edit', 'Edit')}
                        </Menu.Item>

                        <Menu.Item
                            leftSection={<TbTags size={18} />}
                            onClick={() => {
                                showModal('editTagsModal', {
                                    editTagsFrom: 'srsList',
                                    tags: srsList.tags ?? [],
                                    uuid: srsList.uuid
                                })
                            }}
                        >
                            {t('common.field.tags')}
                        </Menu.Item>

                        <Menu.Item
                            color="cyan"
                            leftSection={<PiPulse size={18} />}
                            onClick={() => onCheck(srsList.uuid)}
                        >
                            {tr('srs-lists.feature.check', 'Check availability')}
                        </Menu.Item>

                        <Menu.Item
                            color={srsList.isEnabled ? 'orange' : 'teal'}
                            leftSection={
                                srsList.isEnabled ? (
                                    <PiProhibit size={18} />
                                ) : (
                                    <PiCheck size={18} />
                                )
                            }
                            onClick={() => onToggleEnable(srsList.uuid, !srsList.isEnabled)}
                        >
                            {srsList.isEnabled
                                ? tr('common.disable', 'Disable')
                                : tr('common.enable', 'Enable')}
                        </Menu.Item>

                        <CopyButton timeout={2000} value={srsList.url}>
                            {({ copied, copy }) => (
                                <Menu.Item
                                    color={copied ? 'teal' : undefined}
                                    leftSection={
                                        copied ? <PiCheck size={18} /> : <TbLink size={18} />
                                    }
                                    onClick={copy}
                                >
                                    {copied
                                        ? tr('common.copied', 'Copied')
                                        : tr('srs-lists.feature.copy-url', 'Copy URL')}
                                </Menu.Item>
                            )}
                        </CopyButton>

                        <CopyButton timeout={2000} value={srsList.uuid}>
                            {({ copied, copy }) => (
                                <Menu.Item
                                    color={copied ? 'teal' : undefined}
                                    leftSection={
                                        copied ? <PiCheck size={18} /> : <PiCopy size={18} />
                                    }
                                    onClick={copy}
                                >
                                    {copied
                                        ? tr('common.copied', 'Copied')
                                        : tr('common.action.copy-uuid', 'Copy UUID')}
                                </Menu.Item>
                            )}
                        </CopyButton>

                        <Menu.Item
                            color="red"
                            leftSection={<PiTrashDuotone size={18} />}
                            onClick={() => onDelete(srsList.uuid, displayName)}
                        >
                            {tr('common.delete', 'Delete')}
                        </Menu.Item>
                    </EntityCardShared.Menu>
                </EntityCardShared.Actions>
            </EntityCardShared.Root>
        </WithDndSortable>
    )
}
