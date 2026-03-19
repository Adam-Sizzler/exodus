import { GetSubscriptionInfoByShortUuidCommand } from '@cerberus/backend-contract'

export interface IState {
    subscription: GetSubscriptionInfoByShortUuidCommand.Response['response'] | null
}
