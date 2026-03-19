import {
    TSubscriptionPageLanguageCode,
    TSubscriptionPageRawConfig
} from '@cerberus/subscription-page-types'

export interface IState {
    config: null | TSubscriptionPageRawConfig
    currentLang: TSubscriptionPageLanguageCode
    isConfigLoaded: boolean
}
