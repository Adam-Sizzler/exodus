import {
    ActionIcon,
    Box,
    Button,
    Checkbox,
    Code,
    Group,
    ScrollArea,
    Stack,
    Text,
    Textarea,
    ThemeIcon
} from '@mantine/core'
import {
    TbAlertTriangle,
    TbArrowBackUp,
    TbLock,
    TbLockOpen,
    TbRefresh,
    TbSend,
    TbServer2
} from 'react-icons/tb'
import ReactCountryFlag from 'react-country-flag'
import { useCallback, useState } from 'react'

import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { NodeResponse, useNodePluginExecutor } from '@shared/api/hooks'
import { ActionCardShared } from '@shared/ui/action-card'
import { SectionCard } from '@shared/ui/section-card'

interface IProps {
    nodes: NodeResponse[]
    onClose: () => void
}

type CommandType = 'blockIps' | 'recreateTables' | 'unblockIps'
type Step = 'command' | 'configure' | 'target'

const BLOCK_PLACEHOLDER = `192.168.1.1;0
10.0.0.1;3600
172.16.0.1;60`

const UNBLOCK_PLACEHOLDER = `192.168.1.1
10.0.0.1
172.16.0.1`

const isValidIP = (value: string) => {
    const ipv4 = /^(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}$/
    const ipv6 = /^(([0-9a-f]{1,4}:){2,7}[0-9a-f]{1,4}|::1|::)$/i
    return ipv4.test(value) || ipv6.test(value)
}

export const NodePluginExecutorContent = (props: IProps) => {
    const { nodes, onClose } = props

    const { mutate: executeNodePlugin, isPending } = useNodePluginExecutor({
        mutationFns: {
            onSuccess: () => onClose()
        }
    })

    const [step, setStep] = useState<Step>('command')
    const [command, setCommand] = useState<CommandType | null>(null)
    const [selectedNodeUuids, setSelectedNodeUuids] = useState<Set<string>>(new Set())
    const [blockText, setBlockText] = useState('')
    const [unblockText, setUnblockText] = useState('')
    const [textError, setTextError] = useState<null | string>(null)

    const connectedNodes = nodes.filter((node) => node.isConnected)

    const resetAll = () => {
        setCommand(null)
        setSelectedNodeUuids(new Set())
        setBlockText('')
        setUnblockText('')
        setTextError(null)
    }

    const selectCommand = (cmd: CommandType) => {
        setCommand(cmd)
        setStep(cmd === 'recreateTables' ? 'target' : 'configure')
    }

    const goBack = () => {
        if (step === 'target' && command !== 'recreateTables') {
            setStep('configure')
            return
        }
        setStep('command')
        resetAll()
    }

    const toggleNode = useCallback((uuid: string) => {
        setSelectedNodeUuids((prev) => {
            const next = new Set(prev)
            if (next.has(uuid)) next.delete(uuid)
            else next.add(uuid)
            return next
        })
    }, [])

    const parseBlockText = () => {
        const lines = blockText
            .split('\n')
            .map((line) => line.trim())
            .filter(Boolean)
        const errors: string[] = []
        const entries: { ip: string; timeout: number }[] = []

        lines.forEach((line, index) => {
            const parts = line.split(';')
            const ip = parts[0]?.trim() ?? ''
            const timeout = parseInt(parts[1]?.trim() ?? '0', 10)

            if (!isValidIP(ip)) {
                errors.push(`Line ${index + 1}: "${ip}" is not a valid IP address`)
            } else {
                entries.push({ ip, timeout: Number.isNaN(timeout) ? 0 : timeout })
            }
        })

        return { entries, errors }
    }

    const parseUnblockText = () => {
        const lines = unblockText
            .split('\n')
            .map((line) => line.trim())
            .filter(Boolean)
        const errors: string[] = []
        const ips: string[] = []

        lines.forEach((line, index) => {
            if (!isValidIP(line)) {
                errors.push(`Line ${index + 1}: "${line}" is not a valid IP address`)
            } else {
                ips.push(line)
            }
        })

        return { errors, ips }
    }

    const validateAndProceed = () => {
        if (command === 'blockIps') {
            const { entries, errors } = parseBlockText()
            if (entries.length === 0 && errors.length === 0) {
                setTextError('Enter at least one IP address')
                return false
            }
            if (errors.length > 0) {
                setTextError(errors.join('\n'))
                return false
            }
        }

        if (command === 'unblockIps') {
            const { errors, ips } = parseUnblockText()
            if (ips.length === 0 && errors.length === 0) {
                setTextError('Enter at least one IP address')
                return false
            }
            if (errors.length > 0) {
                setTextError(errors.join('\n'))
                return false
            }
        }

        setTextError(null)
        return true
    }

    const handleSubmit = () => {
        const targetNodes = {
            target: 'specificNodes' as const,
            nodeUuids: Array.from(selectedNodeUuids)
        }

        if (command === 'blockIps') {
            executeNodePlugin({
                variables: { command: { command: 'blockIps', ips: parseBlockText().entries }, targetNodes }
            })
            return
        }

        if (command === 'unblockIps') {
            executeNodePlugin({
                variables: { command: { command: 'unblockIps', ips: parseUnblockText().ips }, targetNodes }
            })
            return
        }

        executeNodePlugin({
            variables: { command: { command: 'recreateTables' }, targetNodes }
        })
    }

    const countryFlag = (countryCode: null | string) => {
        if (!countryCode || countryCode === 'XX') return <TbServer2 size={14} />
        return (
            <ReactCountryFlag
                countryCode={countryCode}
                style={{ fontSize: '1.1em', borderRadius: '2px' }}
            />
        )
    }

    const STEP_MIN_HEIGHT = 380

    if (step === 'command') {
        return (
            <Box mih={STEP_MIN_HEIGHT}>
                <Stack gap="md">
                    <SectionCard.Root>
                        <SectionCard.Section>
                            <BaseOverlayHeader
                                actionIconProps={{ color: 'orange' }}
                                IconComponent={TbAlertTriangle}
                                iconVariant="light"
                                subtitle="Execute plugin maintenance commands on selected online nodes."
                                title="Warning"
                                titleOrder={5}
                            />
                        </SectionCard.Section>
                    </SectionCard.Root>

                    <Stack gap="xs">
                        <ActionCardShared
                            description="Block specific IP addresses on selected nodes"
                            icon={<TbLock size={20} />}
                            onClick={() => selectCommand('blockIps')}
                            title="Block IPs"
                            variant="light"
                        />
                        <ActionCardShared
                            description="Remove IP blocks on selected nodes"
                            icon={<TbLockOpen size={20} />}
                            onClick={() => selectCommand('unblockIps')}
                            title="Unblock IPs"
                            variant="light"
                        />
                        <ActionCardShared
                            description="Recreate plugin runtime rules on selected nodes"
                            icon={<TbRefresh size={20} />}
                            onClick={() => selectCommand('recreateTables')}
                            title="Recreate Tables"
                            variant="light"
                        />
                    </Stack>
                </Stack>
            </Box>
        )
    }

    if (step === 'configure') {
        const isBlock = command === 'blockIps'

        return (
            <Box mih={STEP_MIN_HEIGHT}>
                <SectionCard.Root>
                    <SectionCard.Section>
                        <Group align="flex-start" justify="space-between">
                            <BaseOverlayHeader
                                actionIconProps={{ color: isBlock ? 'cyan' : 'teal' }}
                                IconComponent={isBlock ? TbLock : TbLockOpen}
                                iconVariant="light"
                                subtitle={isBlock ? 'Format: IP;timeout' : 'Format: IP'}
                                title={isBlock ? 'IPs to block' : 'IPs to unblock'}
                                titleOrder={5}
                            />
                            <ActionIcon onClick={goBack} size="lg" variant="default">
                                <TbArrowBackUp size={20} />
                            </ActionIcon>
                        </Group>
                    </SectionCard.Section>

                    <SectionCard.Section>
                        <Text c="dimmed" mb="xs" size="xs">
                            Format one per line: <Code>{isBlock ? 'IP;timeout' : 'IP'}</Code>
                        </Text>
                        <Textarea
                            autosize
                            error={textError}
                            maxRows={10}
                            minRows={5}
                            onChange={(event) => {
                                if (isBlock) setBlockText(event.currentTarget.value)
                                else setUnblockText(event.currentTarget.value)
                                if (textError) setTextError(null)
                            }}
                            placeholder={isBlock ? BLOCK_PLACEHOLDER : UNBLOCK_PLACEHOLDER}
                            styles={{
                                input: {
                                    fontFamily: 'var(--mantine-font-family-monospace)'
                                }
                            }}
                            value={isBlock ? blockText : unblockText}
                        />
                    </SectionCard.Section>
                    <SectionCard.Section>
                        <Group justify="flex-end">
                            <Button
                                onClick={() => {
                                    if (validateAndProceed()) setStep('target')
                                }}
                            >
                                Next
                            </Button>
                        </Group>
                    </SectionCard.Section>
                </SectionCard.Root>
            </Box>
        )
    }

    return (
        <Box mih={STEP_MIN_HEIGHT}>
            <SectionCard.Root>
                <SectionCard.Section>
                    <Group align="flex-start" justify="space-between">
                        <BaseOverlayHeader
                            actionIconProps={{ color: 'violet' }}
                            IconComponent={TbServer2}
                            iconVariant="light"
                            subtitle={`${selectedNodeUuids.size} selected`}
                            title="Nodes"
                            titleOrder={5}
                        />
                        <ActionIcon onClick={goBack} size="lg" variant="default">
                            <TbArrowBackUp size={20} />
                        </ActionIcon>
                    </Group>
                </SectionCard.Section>

                {connectedNodes.length > 0 && (
                    <SectionCard.Section>
                        <ScrollArea.Autosize mah={280} offsetScrollbars>
                            <Stack gap={6}>
                                {connectedNodes.map((node) => {
                                    const isSelected = selectedNodeUuids.has(node.uuid)
                                    return (
                                        <Checkbox.Card
                                            checked={isSelected}
                                            key={node.uuid}
                                            onClick={() => toggleNode(node.uuid)}
                                            p="sm"
                                            radius="md"
                                            style={{
                                                border: isSelected
                                                    ? '1px solid var(--mantine-color-cyan-6)'
                                                    : '1px solid rgba(255,255,255,0.06)',
                                                background: isSelected
                                                    ? 'var(--mantine-color-cyan-light)'
                                                    : 'transparent',
                                                transition: 'all 0.15s ease'
                                            }}
                                        >
                                            <Group gap="sm" wrap="nowrap">
                                                <Checkbox.Indicator size="sm" />
                                                <ThemeIcon
                                                    color={isSelected ? 'cyan' : 'gray'}
                                                    radius="md"
                                                    size="md"
                                                    variant="light"
                                                >
                                                    {countryFlag(node.countryCode)}
                                                </ThemeIcon>
                                                <Stack gap={0} style={{ flex: 1, minWidth: 0 }}>
                                                    <Text fw={600} lineClamp={1} size="sm">
                                                        {node.name}
                                                    </Text>
                                                    {node.address && (
                                                        <Text
                                                            c="dimmed"
                                                            ff="monospace"
                                                            lineClamp={1}
                                                            size="xs"
                                                        >
                                                            {node.address}
                                                        </Text>
                                                    )}
                                                </Stack>
                                            </Group>
                                        </Checkbox.Card>
                                    )
                                })}
                            </Stack>
                        </ScrollArea.Autosize>
                    </SectionCard.Section>
                )}

                {connectedNodes.length === 0 && (
                    <SectionCard.Section>
                        <Text c="dimmed" py="md" size="sm" ta="center">
                            No connected nodes available
                        </Text>
                    </SectionCard.Section>
                )}

                <SectionCard.Section>
                    <Group justify="flex-end">
                        {connectedNodes.length > 0 && (
                            <Button
                                onClick={() => {
                                    if (selectedNodeUuids.size === connectedNodes.length) {
                                        setSelectedNodeUuids(new Set())
                                    } else {
                                        setSelectedNodeUuids(
                                            new Set(connectedNodes.map((node) => node.uuid))
                                        )
                                    }
                                }}
                                size="compact-xs"
                                variant="subtle"
                            >
                                {selectedNodeUuids.size === connectedNodes.length
                                    ? 'Deselect all'
                                    : 'Select all'}
                            </Button>
                        )}

                        <Button
                            color="cyan"
                            disabled={selectedNodeUuids.size === 0}
                            loading={isPending}
                            onClick={handleSubmit}
                            rightSection={<TbSend size={16} />}
                        >
                            Execute
                        </Button>
                    </Group>
                </SectionCard.Section>
            </SectionCard.Root>
        </Box>
    )
}
