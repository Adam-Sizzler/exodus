import { SubscriptionConnectionResponse } from '@shared/api/hooks'

export interface IProps {
    isLoading: boolean
    nodes: SubscriptionConnectionResponse[] | undefined
}
