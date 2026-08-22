import { Anchor, Button, JsonInput, px, Stack, Text } from '@mantine/core'
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

import { BASIC_MUX_PARAMS } from '@shared/constants'

export const MUX_MODAL_ID = 'mux-modal'

interface IProps {
    form: UseFormReturnType<
        | CreateHostCommand.RequestBody
        | UpdateHostCommand.RequestBody
        | UpdateManyHostsCommand.RequestBody
    >
}

export const MuxModalContent = ({ form }: IProps) => {
    const { t } = useTranslation()

    const inputProps = form.getInputProps('muxParams')

    const [value, setValue] = useState<string>(
        (form.getValues().muxParams as unknown as string) ?? ''
    )

    const handleChange = (next: string) => {
        setValue(next)
        form.setFieldValue('muxParams', next)
    }

    return (
        <Stack gap="md">
            <Stack gap={0}>
                <Text c="dimmed" size="sm">
                    {t('base-host-form.this-will-only-be-used-for-xray-json-output')}
                </Text>
                <Text c="dimmed" size="sm">
                    {t('base-host-form.please-ensure-you-provide-a-valid-json-mux-object')}
                </Text>
                <Text c="dimmed" size="sm">
                    {t('base-host-form.for-more-information-refer-to')}{' '}
                    <Anchor
                        href="https://xtls.github.io/ru/config/outbound.html#muxobject"
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
                    handleChange(BASIC_MUX_PARAMS)
                }}
                variant="soft"
            >
                {t('base-host-form.paste-default-mux-params')}
            </Button>

            <Button onClick={() => modals.close(MUX_MODAL_ID)} variant="soft">
                {t('common.close')}
            </Button>

            <JsonInput
                autosize
                error={inputProps.error}
                formatOnBlur
                minRows={15}
                onBlur={inputProps.onBlur}
                onChange={handleChange}
                placeholder={BASIC_MUX_PARAMS}
                validationError={t('base-host-form.invalid-json')}
                value={value}
            />
        </Stack>
    )
}
