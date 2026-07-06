import { ActionIcon, Avatar, Group, Text } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { TbCreditCard, TbTrash } from 'react-icons/tb'

import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'
import { faviconResolver, formatCurrencyWithIntl } from '@shared/utils/misc'
import { formatTimeUtil } from '@shared/utils/time-utils'

import { BillingRecord } from './group-records-by-month'

interface RecordCardProps {
    onDelete: (uuid: string) => void
    record: BillingRecord
}

export function RecordCard(props: RecordCardProps) {
    const { record, onDelete } = props
    const { i18n } = useTranslation()

    return (
        <SectionCard.Root>
            <SectionCard.Section>
                <Group justify="space-between" wrap="nowrap">
                    <BaseOverlayHeader
                        icon={
                            <Avatar
                                alt={record.provider.name}
                                color="initials"
                                imageProps={{ decoding: 'async', loading: 'lazy' }}
                                name={record.provider.name}
                                onLoad={(event) => {
                                    const img = event.target as HTMLImageElement
                                    if (img.naturalWidth <= 24 && img.naturalHeight <= 24) {
                                        img.src = ''
                                    }
                                }}
                                radius="sm"
                                size={18}
                                src={faviconResolver(record.provider.faviconLink)}
                            />
                        }
                        iconColor="teal"
                        IconComponent={TbCreditCard}
                        iconVariant="soft"
                        subtitle={formatTimeUtil({
                            language: i18n.language,
                            template: 'FULL_DATE',
                            time: record.billedAt,
                            utc: true
                        })}
                        title={record.provider.name}
                        truncateTitle
                    />

                    <Group gap="xs" style={{ flexShrink: 0 }} wrap="nowrap">
                        <Text c="teal" fw={700} size="sm" style={{ whiteSpace: 'nowrap' }}>
                            {formatCurrencyWithIntl(record.amount)}
                        </Text>
                        <ActionIcon
                            color="red"
                            onClick={() => onDelete(record.uuid)}
                            size="input-xs"
                            variant="soft"
                        >
                            <TbTrash size={16} />
                        </ActionIcon>
                    </Group>
                </Group>
            </SectionCard.Section>
        </SectionCard.Root>
    )
}
