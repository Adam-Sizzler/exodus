import { SubscriptionConnectionResponse } from '@shared/api/hooks'

export interface IProps {
    fetchedNode?: SubscriptionConnectionResponse | undefined
    node: SubscriptionConnectionResponse
    style?: React.CSSProperties
    withText?: boolean
}
