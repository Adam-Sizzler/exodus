import { NodeResponse } from '@shared/api/hooks'

export interface IProps {
    handleViewNode: (nodeUuid: string) => void
    isDragOverlay?: boolean
    isMobile: boolean
    node: NodeResponse
}
