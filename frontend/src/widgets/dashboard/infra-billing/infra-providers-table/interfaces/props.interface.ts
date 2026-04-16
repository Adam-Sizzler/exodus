import { GetInfraProvidersCommand } from '@exodus/backend-contract'

export interface IProps {
    infraProviders: GetInfraProvidersCommand.Response['response']['providers']
    infraProvidersLoading: boolean
}
