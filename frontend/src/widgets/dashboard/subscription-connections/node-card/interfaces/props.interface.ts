import { SubscriptionConnectionResponse } from '@shared/api/hooks'

export interface IProps {
    handleViewNode: (nodeUuid: string) => void
    index?: number
    isDragOverlay?: boolean
    isMobile: boolean
    node: SubscriptionConnectionResponse
}
