import { GetSubscriptionSettingsCommand } from '@cerberus/backend-contract'

export interface IProps {
    subscriptionSettings: GetSubscriptionSettingsCommand.Response['response']
}
