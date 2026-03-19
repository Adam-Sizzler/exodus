import { GetSubscriptionPageConfigsCommand } from '@cerberus/backend-contract'

export interface IProps {
    configs: GetSubscriptionPageConfigsCommand.Response['response']['configs']
}
