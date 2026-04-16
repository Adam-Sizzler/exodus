import { GetSubscriptionPageConfigsCommand } from '@exodus/backend-contract'

export interface IProps {
    configs: GetSubscriptionPageConfigsCommand.Response['response']['configs']
}
