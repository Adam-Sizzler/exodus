import { GetConfigProfilesCommand } from '@exodus/backend-contract'

import { CONFIG_PROFILES_VIEW_MODE } from '@pages/dashboard/config-profiles/components/interfaces'

export interface IProps {
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles'] | undefined
    setViewMode?: (viewMode: CONFIG_PROFILES_VIEW_MODE) => void
    viewMode?: CONFIG_PROFILES_VIEW_MODE
}
