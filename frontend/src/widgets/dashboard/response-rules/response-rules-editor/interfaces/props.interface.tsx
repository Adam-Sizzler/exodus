import {
    GetSubscriptionSettingsCommand,
    TSubscriptionTemplateType
} from '@cerberus/backend-contract'

export interface IProps {
    groupedTemplates: Record<TSubscriptionTemplateType, string[]>
    responseRules: GetSubscriptionSettingsCommand.Response['response']['responseRules']
    subscriptionSettingsUuid: string
}
