import dayjs from 'dayjs'
import utc from 'dayjs/plugin/utc'

dayjs.extend(utc)

enum ETemplatePreset {
    FULL_DATE = 'D MMMM YYYY',
    FULL_DATETIME = 'D MMMM YYYY HH:mm:ss',
    NUMERIC_DATETIME = 'DD.MM.YYYY HH:mm:ss',
    SHORT_DATE = 'D MMM',
    TIME_FIRST_DATETIME = 'HH:mm:ss, D MMMM YYYY'
}

interface FormatTimeUtilProps {
    language?: string
    template: keyof typeof ETemplatePreset
    time: Date | null | number | string | undefined
    utc?: boolean
}

/**
 * @example
 * formatTimeUtil({ time: date, template: 'SHORT_DATE' })        // "5 Mar"
 * formatTimeUtil({ time: date, template: 'FULL_DATE' })         // "5 March 2026"
 * formatTimeUtil({ time: date, template: 'FULL_DATETIME' })     // "5 March 2026 14:30:00"
 * formatTimeUtil({ time: date, template: 'NUMERIC_DATETIME' })  // "05.03.2026 14:30:00"
 * formatTimeUtil({ time: date, template: 'TIME_FIRST_DATETIME'})// "14:30:00, 5 March 2026"
 */
export const formatTimeUtil = (props: FormatTimeUtilProps): string => {
    const { time, template, language, utc } = props

    if (!time) return '-'

    const date = utc ? dayjs.utc(time) : dayjs(time)
    if (!date.isValid()) return 'Unknown date'

    return (language ? date.locale(language) : date).format(ETemplatePreset[template])
}
