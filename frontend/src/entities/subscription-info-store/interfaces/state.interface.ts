import { GetSubscriptionInfoByShortUuidCommand } from '@exodus/backend-contract'

export interface IState {
    subscription: GetSubscriptionInfoByShortUuidCommand.Response['response'] | null
}
