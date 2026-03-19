import { GetInfraProvidersCommand } from '@cerberus/backend-contract'

export interface IProps {
    infraProviders: GetInfraProvidersCommand.Response['response']['providers']
    infraProvidersLoading: boolean
}
