import { Anchor, Button, JsonInput, px, Select, Stack, Text, Textarea } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { modals } from '@mantine/modals'
import {
    CreateHostCommand,
    UpdateHostCommand,
    UpdateManyHostsCommand
} from '@exodus/backend-contract'
import { useEffect, useState } from 'react'
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
        CreateHostCommand.Request | UpdateHostCommand.Request | UpdateManyHostsCommand.Request
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

// Clash/Mihomo's "smux" block is native YAML (see BASIC_CLASH_MUX_PARAMS),
// unlike Xray/Singbox whose mux params are genuinely JSON. The backend
// (generator_mihomo.go's parseJSONMapString) already accepts this: it first
// tries to parse the stored value as a JSON object, and if that fails,
// falls back to treating it as a JSON-encoded *string* containing YAML text
// (i.e. the raw YAML wrapped/escaped via JSON.stringify). These two helpers
// keep that wrapping entirely on the frontend so the admin only ever sees
// and edits plain, human-readable YAML - the backend contract is untouched.
const wrapClashYamlForStorage = (yamlText: string): string => {
    if (yamlText.trim() === '') {
        return ''
    }
    return JSON.stringify(yamlText)
}

const unwrapClashYamlFromStorage = (stored: string): string => {
    if (stored.trim() === '') {
        return ''
    }
    try {
        const parsed = JSON.parse(stored)
        if (typeof parsed === 'string') {
            return parsed
        }
    } catch {
        // Not a JSON-string-wrapped payload (e.g. legacy/hand-edited value) -
        // fall through and show it as-is, as plain YAML.
    }
    return stored
}

export const MuxModalContent = ({ form }: IProps) => {
    const { t } = useTranslation()

    const [activeCore, setActiveCore] = useState<MuxCore>('SINGBOX')
    const activeField = MUX_FIELD_BY_CORE[activeCore]
    const activePlaceholder = MUX_PLACEHOLDER_BY_CORE[activeCore]
    const isYamlCore = activeCore === 'CLASH'
    const inputProps = form.getInputProps(activeField as never)

    const [value, setValue] = useState<string>(() => {
        const stored =
            ((form.getValues() as Record<string, unknown>)[activeField] as string | undefined) ??
            ''
        return isYamlCore ? unwrapClashYamlFromStorage(stored) : stored
    })

    useEffect(() => {
        const stored =
            ((form.getValues() as Record<string, unknown>)[activeField] as string | undefined) ??
            ''
        setValue(isYamlCore ? unwrapClashYamlFromStorage(stored) : stored)
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [activeField])

    const handleChange = (next: string) => {
        setValue(next)
        const toStore = isYamlCore ? wrapClashYamlForStorage(next) : next
        form.setFieldValue(activeField as never, toStore as never)
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
