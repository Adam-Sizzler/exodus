import { Button, JsonInput, px, Stack, Text } from '@mantine/core'
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

import { BASIC_XHTTP_EXTRA_PARAMS, PASTE_BASIC_XHTTP_EXTRA_PARAMS } from '@shared/constants'

export const XHTTP_MODAL_ID = 'xhttp-modal'

interface IProps {
    form: UseFormReturnType<
        CreateHostCommand.Request | UpdateHostCommand.Request | UpdateManyHostsCommand.Request
    >
}

export const XhttpModalContent = ({ form }: IProps) => {
    const { t } = useTranslation()

    const inputProps = form.getInputProps('xhttpExtraParams')

    const [value, setValue] = useState<string>(
        (form.getValues().xhttpExtraParams as unknown as string) ?? ''
    )

    const handleChange = (next: string) => {
        setValue(next)
        form.setFieldValue('xhttpExtraParams', next)
    }

    return (
        <Stack gap="md">
            <Text c="dimmed" size="sm">
                {t('base-host-form.extra-xhttp-description')}
            </Text>
            <Button
                color="gray"
                leftSection={<TbArrowUp size={px('1.2rem')} />}
                onClick={() => {
                    handleChange(PASTE_BASIC_XHTTP_EXTRA_PARAMS)
                }}
                variant="soft"
            >
                {t('base-host-form.fill-with-sample-xhttp-extra-params')}
            </Button>

            <Button onClick={() => modals.close(XHTTP_MODAL_ID)} variant="soft">
                {t('common.close')}
            </Button>
            <JsonInput
                autosize
                error={inputProps.error}
                formatOnBlur
                minRows={15}
                onBlur={inputProps.onBlur}
                onChange={handleChange}
                placeholder={BASIC_XHTTP_EXTRA_PARAMS}
                validationError={t('base-host-form.invalid-json')}
                value={value}
            />
        </Stack>
    )
}
