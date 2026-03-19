import {
    GetSubscriptionTemplatesCommand,
    TSubscriptionTemplateType
} from '@cerberus/backend-contract'

export interface IProps {
    templates: GetSubscriptionTemplatesCommand.Response['response']['templates']
    templateTitle: string
    type: TSubscriptionTemplateType
}
