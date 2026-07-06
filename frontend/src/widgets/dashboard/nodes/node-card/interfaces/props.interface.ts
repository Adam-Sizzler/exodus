import { GetAllNodesCommand } from '@exodus/backend-contract'

export interface IProps {
    disableReordering?: boolean
    handleViewNode: (nodeUuid: string) => void
    isDragOverlay?: boolean
    isMobile: boolean
    node: GetAllNodesCommand.Response['response'][number]
}
