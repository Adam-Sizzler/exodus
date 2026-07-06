import { UseListStateHandlers } from '@mantine/hooks'
import { GetAllHostsCommand, GetConfigProfilesCommand } from '@exodus/backend-contract'

export interface IProps {
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles'] | undefined
    handlers: UseListStateHandlers<GetAllHostsCommand.Response['response'][number]>
    hosts: GetAllHostsCommand.Response['response'] | undefined
    selectedHosts: string[]
    setSelectedHosts: React.Dispatch<React.SetStateAction<string[]>>
    state: GetAllHostsCommand.Response['response']
}
