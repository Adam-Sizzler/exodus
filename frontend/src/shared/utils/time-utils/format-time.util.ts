import dayjs from 'dayjs'

enum ETemplatePreset {
    FULL_DATE = 'D MMMM YYYY',
    FULL_DATETIME = 'D MMMM YYYY HH:mm:ss',
    NUMERIC_DATETIME = 'DD.MM.YYYY HH:mm:ss',
    SHORT_DATE = 'D MMM',
    TIME_FIRST_DATETIME = 'HH:mm:ss, D MMMM YYYY'
}

type TLegacyTemplatePreset = 'DD.MM.YYYY HH:mm:ss' | 'D MMM' | 'D MMMM YYYY'

interface FormatTimeUtilProps {
    language?: string
    template: keyof typeof ETemplatePreset
    time: Date | null | number | string | undefined
}

export function formatTimeUtil(props: FormatTimeUtilProps): string
export function formatTimeUtil(
    time: null | number | string | undefined,
    template: TLegacyTemplatePreset
): string
export function formatTimeUtil(
    propsOrTime: FormatTimeUtilProps | null | number | string | undefined,
    legacyTemplate?: TLegacyTemplatePreset
): string {
    const isObjectCall =
        typeof propsOrTime === 'object' && propsOrTime !== null && 'time' in propsOrTime
    const time = isObjectCall ? propsOrTime.time : propsOrTime
    const template = isObjectCall ? ETemplatePreset[propsOrTime.template] : legacyTemplate
    const language = isObjectCall ? propsOrTime.language : undefined

    if (!time) return 'Unknown date'

    const date = dayjs(time)
    if (!date.isValid()) return 'Unknown date'

    return (language ? date.locale(language) : date).format(template)
}
