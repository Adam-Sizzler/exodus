import type { InputBaseProps } from '@mantine/core'

import { GetConfigProfilesCommand } from '@exodus/backend-contract'

export interface IProps extends InputBaseProps {
    nodes: GetConfigProfilesCommand.Response['response']['configProfiles'][number]['nodes']
}
