import { modals } from '@mantine/modals'
import { useTranslation } from 'react-i18next'

import { useDeleteInfraBillingHistoryRecord } from '@shared/api/hooks'

export function useDeleteBillingRecord(refetchRecords: () => void) {
    const { t } = useTranslation()

    const { mutate: deleteRecord } = useDeleteInfraBillingHistoryRecord({
        mutationFns: {
            onSuccess: () => {
                refetchRecords()
            }
        }
    })

    return (uuid: string) =>
        modals.openConfirmModal({
            title: t('common.confirm-action'),
            children: t('common.confirm-action-description'),
            labels: { confirm: t('common.delete'), cancel: t('common.cancel') },
            centered: true,
            confirmProps: { color: 'red' },
            onConfirm: () => deleteRecord({ route: { uuid } })
        })
}
