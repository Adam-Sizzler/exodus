import { EntityCardActions } from './entity-card.actions'
import { EntityCardButtonAction } from './entity-card.button-action'
import { EntityCardContent } from './entity-card.content'
import { EntityCardHeader } from './entity-card.header'
import { EntityCardIcon } from './entity-card.icon'
import { EntityCardMenuAction } from './entity-card.menu-action'
import { EntityCardRoot } from './entity-card.root'

export const EntityCardShared = {
    Root: EntityCardRoot,
    Actions: EntityCardActions,
    Button: EntityCardButtonAction,
    Content: EntityCardContent,
    Header: EntityCardHeader,
    Icon: EntityCardIcon,
    Menu: EntityCardMenuAction
}
