import { SubscriptionConnectionResponse } from '@shared/api/hooks'

export interface IProps {
    handleClose: () => void
    node: SubscriptionConnectionResponse
}
