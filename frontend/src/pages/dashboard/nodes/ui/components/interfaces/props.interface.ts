import { NodeResponse } from '@shared/api/hooks'

export interface IProps {
    isLoading: boolean
    nodes: NodeResponse[] | undefined
}
