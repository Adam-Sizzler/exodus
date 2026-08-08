import { Anchor, Button, JsonInput, px, Select, Stack, Text, Textarea } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { modals } from '@mantine/modals'
import {
    CreateHostCommand,
    UpdateHostCommand,
    UpdateManyHostsCommand
} from '@exodus/backend-contract'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { TbArrowUp } from 'react-icons/tb'

import {
    BASIC_CLASH_MUX_PARAMS,
    BASIC_SINGBOX_MUX_PARAMS,
    BASIC_XRAY_MUX_PARAMS
} from '@shared/constants'

export const MUX_MODAL_ID = 'mux-modal'

interface IProps {
    form: UseFormReturnType<
        | CreateHostCommand.RequestBody
        | UpdateHostCommand.RequestBody
        | UpdateManyHostsCommand.RequestBody
    >
}

const MUX_CORE_OPTIONS = [
    { value: 'SINGBOX', label: 'Singbox' },
    { value: 'CLASH', label: 'Clash/Mihomo' },
    { value: 'XRAY', label: 'Xray' }
] as const

type MuxCore = (typeof MUX_CORE_OPTIONS)[number]['value']

const MUX_FIELD_BY_CORE: Record<MuxCore, 'clashMuxParams' | 'muxParams' | 'singboxMuxParams'> = {
    CLASH: 'clashMuxParams',
    SINGBOX: 'singboxMuxParams',
    XRAY: 'muxParams'
}

const MUX_PLACEHOLDER_BY_CORE: Record<MuxCore, string> = {
    CLASH: BASIC_CLASH_MUX_PARAMS,
    SINGBOX: BASIC_SINGBOX_MUX_PARAMS,
    XRAY: BASIC_XRAY_MUX_PARAMS
}

const MUX_DOCS_BY_CORE: Record<MuxCore, string> = {
    CLASH: 'https://wiki.metacubex.one/en/config/proxies/',
    SINGBOX: 'https://sing-box.sagernet.org/configuration/shared/multiplex/',
    XRAY: 'https://xtls.github.io/ru/config/outbound.html#muxobject'
}

export const MuxModalContent = ({ form }: IProps) => {
    const { t } = useTranslation()

    const [activeCore, setActiveCore] = useState<MuxCore>('SINGBOX')
    const activeField = MUX_FIELD_BY_CORE[activeCore]
    const activePlaceholder = MUX_PLACEHOLDER_BY_CORE[activeCore]
    const isYamlCore = activeCore === 'CLASH'
    const inputProps = form.getInputProps(activeField as never)

    const [value, setValue] = useState<string>(
        () =>
            ((form.getValues() as Record<string, unknown>)[activeField] as string | undefined) ?? ''
    )

    // Reset `value` synchronously, during render, whenever the active core
    // (and therefore the underlying form field) changes - not in a
    // useEffect. useEffect runs *after* the DOM has already been painted
    // with the previous core's raw value, and Mantine's JsonInput only
    // re-validates on its own onChange, not when its `value` prop changes
    // from outside - so for one visible frame it would show "Invalid JSON"
    // for the old core's text (e.g. Clash's YAML) rendered inside the new
    // core's JSON validator, only clearing once the user typed or clicked.
    // Adjusting state directly during render (React's documented pattern
    // for "resetting state when a prop changes") re-renders before the
    // browser paints anything, so that stale frame never appears.
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
                data={MUX_CORE_OPTIONS.map((option) => ({
                    label: option.label,
                    value: option.value
                }))}
                label={t('base-host-form.mux-core')}
                onChange={(next) => setActiveCore((next as MuxCore) || 'SINGBOX')}
                value={activeCore}
            />

            <Stack gap={0}>
                <Text c="dimmed" size="sm">
                    {t('base-host-form.this-will-only-be-used-for-selected-core-output', {
                        core:
                            MUX_CORE_OPTIONS.find((option) => option.value === activeCore)?.label ??
                            'Singbox'
                    })}
                </Text>
                <Text c="dimmed" size="sm">
                    {isYamlCore
                        ? t('base-host-form.please-ensure-you-provide-a-valid-yaml-mux-object')
                        : t('base-host-form.please-ensure-you-provide-a-valid-json-mux-object')}
                </Text>
                <Text c="dimmed" size="sm">
                    {t('base-host-form.for-more-information-refer-to')}{' '}
                    <Anchor
                        href={MUX_DOCS_BY_CORE[activeCore]}
                        rel="noopener noreferrer"
                        target="_blank"
                    >
                        {t('base-host-form.xtls-documentation')}
                    </Anchor>
                    .
                </Text>
            </Stack>
            <Button
                color="gray"
                leftSection={<TbArrowUp size={px('1.2rem')} />}
                onClick={() => {
                    handleChange(activePlaceholder)
                }}
                variant="soft"
            >
                {t('base-host-form.paste-default-mux-params')}
            </Button>

            <Button onClick={() => modals.close(MUX_MODAL_ID)} variant="soft">
                {t('common.close')}
            </Button>

            {isYamlCore ? (
                <Textarea
                    autosize
                    error={inputProps.error}
                    minRows={15}
                    onBlur={inputProps.onBlur}
                    onChange={(event) => handleChange(event.currentTarget.value)}
                    placeholder={activePlaceholder}
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
                    placeholder={activePlaceholder}
                    validationError={t('base-host-form.invalid-json')}
                    value={value}
                />
            )}
        </Stack>
    )
}
