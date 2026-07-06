import { ActionIcon, Tooltip } from '@mantine/core'
import { modals } from '@mantine/modals'
import { UserHwidDevicesContentWidget } from '@widgets/dashboard/users/user-hwid-devices-drawer/user-hwid-devices.content.widget'
import { useTranslation } from 'react-i18next'
import { TbDevices } from 'react-icons/tb'

import { useIsMobile } from '@shared/hooks'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { IProps } from './interfaces'

export function GetHwidUserDevicesFeature(props: IProps) {
    const { userUuid } = props
    const { t } = useTranslation()

    const isMobile = useIsMobile()

    return (
        <Tooltip label={t('get-hwid-user-devices.feature.hwid-devices')}>
            <ActionIcon
                color="indigo"
                onClick={() => {
                    modals.open({
                        title: (
                            <BaseOverlayHeader
                                iconColor="violet"
                                IconComponent={TbDevices}
                                iconVariant="soft"
                                title={t('get-hwid-user-devices.feature.hwid-devices')}
                            />
                        ),
                        centered: true,
                        size: 'min(1500px, 90vw)',
                        fullScreen: isMobile,
                        styles: {
                            body: {
                                height: 'calc(100% - 60px)',
                                display: 'flex',
                                flexDirection: 'column'
                            }
                        },
                        children: (
                            <UserHwidDevicesContentWidget mobile={isMobile} userUuid={userUuid} />
                        )
                    })
                }}
                size="lg"
                variant="soft"
            >
                <TbDevices size="22px" />
            </ActionIcon>
        </Tooltip>
    )
}
