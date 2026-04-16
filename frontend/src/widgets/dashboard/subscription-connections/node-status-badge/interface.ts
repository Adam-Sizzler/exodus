import { GetAllNodesCommand } from '@exodus/backend-contract'

export interface IProps {
    fetchedNode?: GetAllNodesCommand.Response['response'][number] | undefined
    node: GetAllNodesCommand.Response['response'][number]
    style?: React.CSSProperties
    withText?: boolean
}
