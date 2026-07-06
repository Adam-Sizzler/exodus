import { Anchor, Button, JsonInput, Stack, Text } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { modals } from '@mantine/modals'
import {
    CreateHostCommand,
    UpdateHostCommand,
    UpdateManyHostsCommand
} from '@exodus/backend-contract'
import { Trans, useTranslation } from 'react-i18next'
import { TbDeviceFloppy } from 'react-icons/tb'

export const FINAL_MASK_MODAL_ID = 'final-mask-modal'

const FINAL_MASK_PLACEHOLDER = {
    tcp: [
        {
            type: '',
            settings: {}
        }
    ],
    udp: [
        {
            type: '',
            settings: {}
        }
    ],
    quicParams: {}
}
interface IProps {
    form: UseFormReturnType<
        CreateHostCommand.Request | UpdateHostCommand.Request | UpdateManyHostsCommand.Request
    >
}

export const FinalMaskModalContent = ({ form }: IProps) => {
    const { t } = useTranslation()
    return (
        <Stack gap="md">
            <Text c="dimmed" size="sm">
                <Trans
                    components={{
                        anchor: (
                            <Anchor
                                href="https://xtls.github.io/ru/config/transports/finalmask.html"
                                rel="noopener noreferrer"
                                target="_blank"
                            />
                        )
                    }}
                    i18nKey="base-host-form.final-mask-description"
                />
            </Text>
            <Button
                leftSection={<TbDeviceFloppy />}
                onClick={() => modals.close(FINAL_MASK_MODAL_ID)}
                size="md"
                variant="soft"
            >
                {t('common.close')}
            </Button>
            <JsonInput
                autosize
                formatOnBlur
                key={form.key('finalMask')}
                maxRows={20}
                minRows={15}
                placeholder={JSON.stringify(FINAL_MASK_PLACEHOLDER, null, 2)}
                validationError={t('base-host-form.invalid-json')}
                {...form.getInputProps('finalMask')}
            />
        </Stack>
    )
}
