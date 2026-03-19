import LanguageDetector from 'i18next-browser-languagedetector/cjs'
import { initReactI18next } from 'node_modules/react-i18next'
import HttpApi from 'i18next-http-backend'
import i18n from 'i18next'
import { withBasePath } from '@shared/constants/base-path'

i18n.use(initReactI18next)
    .use(LanguageDetector)
    .use(HttpApi)
    .init({
        fallbackLng: 'en',
        debug: process.env.NODE_ENV === 'development',
        defaultNS: ['cerberus'],
        ns: ['cerberus'],
        detection: {
            order: ['localStorage', 'navigator', 'htmlTag', 'path', 'subdomain'],
            convertDetectedLanguage: (lng) => (lng.includes('-') ? lng.split('-')[0] : lng)
        },
        load: 'languageOnly',
        preload: ['en', 'ru', 'fa', 'zh'],
        backend: {
            loadPath: withBasePath('/locales/{{lng}}/{{ns}}.json')
        },
        interpolation: {
            escapeValue: false
        },
        react: {
            useSuspense: true
        }
    })

export default i18n
