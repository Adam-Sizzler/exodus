import { Button, JsonInput, Select, Stack, Text, Textarea } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { modals } from '@mantine/modals'
import {
    CreateHostCommand,
    UpdateHostCommand,
    UpdateManyHostsCommand
} from '@exodus/backend-contract'
import { useState } from 'react'

export const CUSTOM_PARAMS_MODAL_ID = 'custom-params-modal'

interface IProps {
    form: UseFormReturnType<
        CreateHostCommand.Request | UpdateHostCommand.Request | UpdateManyHostsCommand.Request
    >
}

const CUSTOM_CORE_OPTIONS = [
    { value: 'SINGBOX', label: 'Singbox' },
    { value: 'MIHOMO', label: 'Mihomo' }
] as const

type CustomCore = (typeof CUSTOM_CORE_OPTIONS)[number]['value']

const CUSTOM_FIELD_BY_CORE: Record<CustomCore, 'singboxCustomParams' | 'mihomoCustomParams'> = {
    SINGBOX: 'singboxCustomParams',
    MIHOMO: 'mihomoCustomParams'
}

export const CustomParamsModalContent = ({ form }: IProps) => {
    const [activeCore, setActiveCore] = useState<CustomCore>('SINGBOX')
    const activeField = CUSTOM_FIELD_BY_CORE[activeCore]
    const isYamlCore = activeCore === 'MIHOMO'
    const inputProps = form.getInputProps(activeField as never)

    const [value, setValue] = useState<string>(
        () =>
            ((form.getValues() as Record<string, unknown>)[activeField] as string | undefined) ?? ''
    )

    const [trackedField, setTrackedField] = useState(activeField)
    if (activeField !== trackedField) {
        setTrackedField(activeField)
        setValue(
            ((form.getValues() as Record<string, unknown>)[activeField] as string | undefined) ?? ''
        )
    }

    const handleChange = (next: string) => {
        setValue(next)
        form.setFieldValue(activeField as never, next as never)
    }

    return (
        <Stack gap="md">
            <Select
                allowDeselect={false}
                data={CUSTOM_CORE_OPTIONS.map((option) => ({
                    label: option.label,
                    value: option.value
                }))}
                label="Core"
                onChange={(next) => setActiveCore((next as CustomCore) || 'SINGBOX')}
                value={activeCore}
            />

            <Stack gap={0}>
                <Text c="dimmed" size="sm">
                    Provide valid {isYamlCore ? 'YAML' : 'JSON'} for custom parameters for the {activeCore} core.
                </Text>
            </Stack>

            <Button onClick={() => modals.close(CUSTOM_PARAMS_MODAL_ID)} variant="soft">
                Close
            </Button>

            {isYamlCore ? (
                <Textarea
                    autosize
                    error={inputProps.error}
                    minRows={15}
                    onBlur={inputProps.onBlur}
                    onChange={(event) => handleChange(event.currentTarget.value)}
                    placeholder="packet-encoding: xudp"
                    styles={{ input: { fontFamily: 'monospace' } }}
                    value={value}
                />
            ) : (
                <JsonInput
                    autosize
                    error={inputProps.error}
                    formatOnBlur
                    minRows={15}
                    onBlur={inputProps.onBlur}
                    onChange={handleChange}
                    placeholder={`{\n  "packet_encoding": "xudp"\n}`}
                    validationError="Invalid JSON"
                    value={value}
                />
            )}
        </Stack>
    )
}
