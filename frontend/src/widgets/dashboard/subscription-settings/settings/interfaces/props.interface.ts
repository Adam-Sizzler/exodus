import { GetSubscriptionSettingsCommand } from '@exodus/backend-contract'

export interface IProps {
    subscriptionSettings: GetSubscriptionSettingsCommand.Response['response']
}
