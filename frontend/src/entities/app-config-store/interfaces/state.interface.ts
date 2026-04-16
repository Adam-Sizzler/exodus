import {
    TSubscriptionPageLanguageCode,
    TSubscriptionPageRawConfig
} from '@exodus/subscription-page-types'

export interface IState {
    config: null | TSubscriptionPageRawConfig
    currentLang: TSubscriptionPageLanguageCode
    isConfigLoaded: boolean
}
