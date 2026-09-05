import {
    Button,
    Card,
    Divider,
    Group,
    Modal,
    Select,
    Stack,
    TagsInput,
    Text,
    TextInput,
    Textarea,
    ThemeIcon,
    Title
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { PiInfo, PiListChecks } from 'react-icons/pi'
import { TbPlus } from 'react-icons/tb'

import { SRSListsHeaderActionButtonsFeature } from '@features/ui/dashboard/srs-lists/header-action-buttons'
import { QueryKeys, useCreateSRSLists, useGetSRSListsTags, useUpdateSRSList } from '@shared/api/hooks'
import { queryClient } from '@shared/api/query-client'
import { Page } from '@shared/ui/page'
import { PageHeaderShared } from '@shared/ui/page-header/page-header.shared'
import { TagInputPill } from '@shared/ui/tag-input-pill'
import {
    SRSListItem,
    SRSListsGridWidget,
    SRSListsSpotlightWidget
} from '@widgets/dashboard/srs-lists'

interface Props {
    srsLists: SRSListItem[]
}

export function SRSListsPageComponent(props: Props) {
    const { srsLists } = props
    const { t } = useTranslation()
    const tr = (key: string, defaultValue: string) => t(key, { defaultValue })

    const { data: knownTagsData } = useGetSRSListsTags()

    const [createModalOpened, createModalHandlers] = useDisclosure(false)
    const [editModalOpened, editModalHandlers] = useDisclosure(false)

    const [editingItem, setEditingItem] = useState<null | SRSListItem>(null)
    const [urlsText, setUrlsText] = useState('')
    const [createTags, setCreateTags] = useState<string[]>([])
    const [updateInterval, setUpdateInterval] = useState('1d')
    const [createEnabled, setCreateEnabled] = useState(true)
    const [editURL, setEditURL] = useState('')
    const [editUpdateInterval, setEditUpdateInterval] = useState('1d')
    const [editEnabled, setEditEnabled] = useState(true)

    const refetchSRSLists = async () => {
        await queryClient.refetchQueries({ queryKey: QueryKeys.srsLists.getSRSLists.queryKey })
        await queryClient.invalidateQueries({ queryKey: QueryKeys.srsLists.getSRSListsTags.queryKey })
    }

    const { mutate: createSRSLists, isPending: isCreating } = useCreateSRSLists({
        mutationFns: {
            onSuccess: async () => {
                createModalHandlers.close()
                setUrlsText('')
                setCreateTags([])
                setUpdateInterval('1d')
                setCreateEnabled(true)
                await refetchSRSLists()
            }
        }
    })

    const { mutate: updateSRSList, isPending: isUpdating } = useUpdateSRSList({
        mutationFns: {
            onSuccess: async () => {
                editModalHandlers.close()
                setEditingItem(null)
                await refetchSRSLists()
            }
        }
    })

    const openEditModal = (item: SRSListItem) => {
        setEditingItem(item)
        setEditURL(item.url)
        setEditUpdateInterval(item.updateInterval || '1d')
        setEditEnabled(item.isEnabled)
        editModalHandlers.open()
    }

    const onSubmitEdit = () => {
        if (!editingItem) return
        updateSRSList({
            variables: {
                uuid: editingItem.uuid,
                url: editURL.trim(),
                updateInterval: editUpdateInterval.trim() || '1d',
                isEnabled: editEnabled
            }
        })
    }

    const onSubmitCreate = () => {
        const urls = urlsText
            .split(/\n|,|\s+/g)
            .map((value) => value.trim())
            .filter(Boolean)

        if (urls.length === 0) return

        createSRSLists({
            variables: {
                urls,
                tags: createTags,
                updateInterval: updateInterval.trim() || '1d',
                isEnabled: createEnabled
            }
        })
    }

    return (
        <Page title={tr('constants.srs-lists', 'SRS Lists')}>
            <PageHeaderShared
                actions={
                    <SRSListsHeaderActionButtonsFeature
                        onOpenCreateModal={createModalHandlers.open}
                        srsListUuids={srsLists.map((i) => i.uuid)}
                    />
                }
                icon={<PiListChecks size={24} />}
                title={tr('constants.srs-lists', 'SRS Lists')}
            />

            <SRSListsGridWidget
                onCreateItem={createModalHandlers.open}
                onEditItem={openEditModal}
                srsLists={srsLists}
            />

            {srsLists.length > 0 && (
                <SRSListsSpotlightWidget onEditItem={openEditModal} srsLists={srsLists} />
            )}

            {/* Create SRS Lists Modal */}
            <Modal
                centered
                onClose={createModalHandlers.close}
                opened={createModalOpened}
                size="lg"
                title={
                    <Group gap="sm" wrap="nowrap">
                        <ThemeIcon color="teal" size="lg" variant="soft">
                            <PiListChecks size={20} />
                        </ThemeIcon>
                        <Title order={4}>
                            {tr('srs-lists.feature.add-links', 'Add SRS Lists')}
                        </Title>
                    </Group>
                }
            >
                <Stack gap="md">
                    <Card p="md" withBorder>
                        <Stack gap="xs">
                            <Group gap="sm" wrap="nowrap">
                                <ThemeIcon color="cyan" size="md" variant="soft">
                                    <PiInfo size={16} />
                                </ThemeIcon>
                                <Title order={5}>
                                    {tr(
                                        'srs-lists.feature.direct-link-required',
                                        'Direct link required'
                                    )}
                                </Title>
                            </Group>
                            <Divider opacity={0.3} />
                            <Text c="dimmed" size="xs">
                                {tr(
                                    'srs-lists.feature.direct-url-hint',
                                    'Use a direct .srs file URL (example: https://raw.githubusercontent.com/.../ruleset.srs). Links with /blob/ are web pages and will not work.'
                                )}
                            </Text>
                            <Textarea
                                autosize
                                description={tr(
                                    'srs-lists.feature.multiple-urls-hint',
                                    'One or multiple URLs separated by spaces, commas or new lines'
                                )}
                                label={tr('srs-lists.feature.urls', 'URLs')}
                                minRows={5}
                                onChange={(event) => setUrlsText(event.currentTarget.value)}
                                placeholder="https://raw.githubusercontent.com/.../ruleset.srs"
                                required
                                value={urlsText}
                            />
                        </Stack>
                    </Card>

                    <TagsInput
                        clearable
                        data={knownTagsData?.tags ?? []}
                        description={tr(
                            'srs-lists.feature.tags-hint',
                            'Optional categorization tags for filtering (e.g. PROD, STREAMING)'
                        )}
                        label={t('common.field.tags')}
                        maxTags={10}
                        onChange={(next) => {
                            const normalized = [
                                ...new Set(next.map((tag) => tag.trim().toUpperCase()).filter(Boolean))
                            ]
                            setCreateTags(normalized)
                        }}
                        placeholder="ENV:PROD, STREAMING"
                        renderPill={({ value: tag, onRemove }) => (
                            <TagInputPill onRemove={onRemove} value={tag} />
                        )}
                        splitChars={[',', ' ', ';']}
                        value={createTags}
                    />

                    <TextInput
                        description={tr(
                            'srs-lists.feature.update-interval-hint',
                            'Supported by sing-box duration parser, e.g.: 10m, 1h, 12h, 1d, 7d.'
                        )}
                        label={tr('srs-lists.feature.update-interval', 'Update interval')}
                        onChange={(event) => setUpdateInterval(event.currentTarget.value)}
                        placeholder="1d"
                        value={updateInterval}
                    />

                    <Select
                        data={[
                            {
                                value: 'enabled',
                                label: tr('srs-lists.feature.enabled', 'Enabled')
                            },
                            {
                                value: 'disabled',
                                label: tr('srs-lists.feature.disabled', 'Disabled')
                            }
                        ]}
                        label={tr('srs-lists.feature.status', 'Status')}
                        onChange={(value) => setCreateEnabled(value !== 'disabled')}
                        value={createEnabled ? 'enabled' : 'disabled'}
                    />

                    <Button
                        color="teal"
                        fullWidth
                        leftSection={<TbPlus size={16} />}
                        loading={isCreating}
                        onClick={onSubmitCreate}
                    >
                        {tr('common.save', 'Save')}
                    </Button>
                </Stack>
            </Modal>

            {/* Edit SRS List Modal */}
            <Modal
                centered
                onClose={() => {
                    editModalHandlers.close()
                    setEditingItem(null)
                }}
                opened={editModalOpened}
                size="md"
                title={
                    <Group gap="sm" wrap="nowrap">
                        <ThemeIcon color="blue" size="lg" variant="soft">
                            <PiListChecks size={20} />
                        </ThemeIcon>
                        <Title order={4}>
                            {tr('srs-lists.feature.edit-list', 'Edit SRS list')}
                        </Title>
                    </Group>
                }
            >
                <Stack gap="md">
                    <TextInput
                        label={tr('srs-lists.feature.url', 'URL')}
                        onChange={(event) => setEditURL(event.currentTarget.value)}
                        placeholder="https://.../ruleset.srs"
                        required
                        value={editURL}
                    />
                    <TextInput
                        description={tr(
                            'srs-lists.feature.update-interval-hint',
                            'e.g. 10m, 1h, 12h, 1d, 7d'
                        )}
                        label={tr('srs-lists.feature.update-interval', 'Update interval')}
                        onChange={(event) => setEditUpdateInterval(event.currentTarget.value)}
                        placeholder="1d"
                        value={editUpdateInterval}
                    />
                    <Select
                        data={[
                            {
                                value: 'enabled',
                                label: tr('srs-lists.feature.enabled', 'Enabled')
                            },
                            {
                                value: 'disabled',
                                label: tr('srs-lists.feature.disabled', 'Disabled')
                            }
                        ]}
                        label={tr('srs-lists.feature.status', 'Status')}
                        onChange={(value) => setEditEnabled(value !== 'disabled')}
                        value={editEnabled ? 'enabled' : 'disabled'}
                    />

                    <Button
                        color="teal"
                        fullWidth
                        loading={isUpdating}
                        onClick={onSubmitEdit}
                    >
                        {tr('common.save', 'Save')}
                    </Button>
                </Stack>
            </Modal>
        </Page>
    )
}
