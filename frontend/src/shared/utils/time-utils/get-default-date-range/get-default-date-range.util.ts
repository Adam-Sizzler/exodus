import dayjs from 'dayjs'

interface DefaultDateRange {
    end: string
    start: string
}

/**
 * Builds the default date range (last `days` days, ending today) as
 * `YYYY-MM-DD` strings.
 *
 * Must be called at component mount (e.g. from a lazy `useState` initializer),
 * NOT stored in a module-level `const`. A top-level const evaluates `dayjs()`
 * once when the bundle is first parsed, so a long-lived tab/PWA freezes "today"
 * at the day it was opened — making the range silently drift further into the
 * past the longer the page stays open.
 */
export const getDefaultDateRange = (days = 7): DefaultDateRange => ({
    start: dayjs()
        .subtract(days - 1, 'day')
        .format('YYYY-MM-DD'),
    end: dayjs().format('YYYY-MM-DD')
})
