import { SRSListItem } from '../srs-list-card'

export interface IProps {
    onCreateItem?: () => void
    onEditItem: (item: SRSListItem) => void
    srsLists: SRSListItem[]
}
