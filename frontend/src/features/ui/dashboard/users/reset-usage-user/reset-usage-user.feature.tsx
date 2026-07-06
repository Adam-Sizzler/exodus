import { Loader, Menu } from '@mantine/core'
import { modals } from '@mantine/modals'
import { useTranslation } from 'react-i18next'
import { PiClockCounterClockwiseDuotone } from 'react-icons/pi'

import { queryClient } from '@shared/api'
import { useResetUserTraffic, usersQueryKeys } from '@shared/api/hooks'

import { IProps } from './interfaces'

export function ResetUsageUserFeature(props: IProps) {
    const { userUuid } = props
    const { t } = useTranslation()

    const { mutate: resetUserTraffic, isPending: isResetUserTrafficPending } = useResetUserTraffic({
        mutationFns: {
            onSuccess: (data) => {
                queryClient.setQueryData(
                    usersQueryKeys.getUserByUuid({ uuid: userUuid }).queryKey,
                    data
                )
            }
        }
    })

    const handleResetUsage = async () => {
        resetUserTraffic({
            route: {
                uuid: userUuid
            }
        })
    }

    const openModal = () =>
        modals.openConfirmModal({
            title: t('common.confirm-action'),
            children: t('common.confirm-action-description'),
            labels: { confirm: t('reset-usage-user.feature.reset'), cancel: t('common.cancel') },
            centered: true,
            confirmProps: { color: 'red', variant: 'soft' },
            onConfirm: () => handleResetUsage()
        })

    return (
        <Menu.Item
            leftSection={
                isResetUserTrafficPending ? (
                    <Loader size="1rem" />
                ) : (
                    <PiClockCounterClockwiseDuotone size="16px" />
                )
            }
            onClick={openModal}
        >
            {t('reset-usage-user.feature.reset-usage')}
        </Menu.Item>
    )
}
