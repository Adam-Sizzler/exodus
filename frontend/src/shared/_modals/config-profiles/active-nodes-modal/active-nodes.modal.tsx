import NiceModal, { useModal } from '@ebay/nice-modal-react'
import {
    ActionIcon,
    Card,
    Center,
    CopyButton,
    Group,
    Modal,
    Stack,
    Text,
    Tooltip
} from '@mantine/core'
import { GetConfigProfilesCommand } from '@exodus/backend-contract'
import ReactCountryFlag from 'react-country-flag'
import { useTranslation } from 'react-i18next'
import { PiCheck, PiCopy, PiCpu } from 'react-icons/pi'
import { TbCpu } from 'react-icons/tb'

import { useNiceMantineModal } from '@shared/_modals/use-nice-modal'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

interface IProps {
    nodes: GetConfigProfilesCommand.Response['response']['configProfiles'][number]['nodes']
    profileName: string
}

export const ActiveNodesModal = NiceModal.create((props: IProps) => {
    const { nodes, profileName } = props

    const modal = useModal()
    const { modalProps } = useNiceMantineModal({ modal })

    const { t } = useTranslation()

    return (
        <Modal
            {...modalProps}
            size="lg"
            title={
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={TbCpu}
                    iconVariant="soft"
                    title={profileName}
                />
            }
        >
            <Stack gap="md">
                {nodes.length > 0 ? (
                    <Stack gap="xs">
                        {nodes.map((node) => (
                            <Card key={node.uuid} p="sm" withBorder>
                                <Group align="center" justify="space-between" wrap="nowrap">
                                    <Group
                                        align="center"
                                        gap="sm"
                                        style={{
                                            flex: 1,
                                            minWidth: 0
                                        }}
                                    >
                                        {node.countryCode && node.countryCode !== 'XX' && (
                                            <ReactCountryFlag
                                                countryCode={node.countryCode}
                                                style={{
                                                    fontSize: '1.1em',
                                                    borderRadius: '2px'
                                                }}
                                            />
                                        )}

                                        <Text
                                            fw={500}
                                            size="sm"
                                            style={{
                                                flex: 1
                                            }}
                                            truncate
                                        >
                                            {node.name}
                                        </Text>
                                    </Group>

                                    <CopyButton timeout={2000} value={node.uuid}>
                                        {({ copied, copy }) => (
                                            <Tooltip label={copied ? 'Copied!' : 'Copy UUID'}>
                                                <ActionIcon
                                                    color={copied ? 'teal' : 'gray'}
                                                    onClick={copy}
                                                    size="sm"
                                                    variant="subtle"
                                                >
                                                    {copied ? (
                                                        <PiCheck size={14} />
                                                    ) : (
                                                        <PiCopy size={14} />
                                                    )}
                                                </ActionIcon>
                                            </Tooltip>
                                        )}
                                    </CopyButton>
                                </Group>
                            </Card>
                        ))}
                    </Stack>
                ) : (
                    <Center py="xl">
                        <Stack align="center" gap="sm">
                            <PiCpu
                                size={48}
                                style={{
                                    color: 'var(--mantine-color-gray-5)'
                                }}
                            />
                            <Text c="dimmed" size="sm" ta="center">
                                {t(
                                    'active-nodes-list.modal.shared.this-profile-is-not-active-on-any-nodes'
                                )}
                            </Text>
                        </Stack>
                    </Center>
                )}
            </Stack>
        </Modal>
    )
})
