import { createTheme } from '@mantine/core'

import { variantColorResolver } from './colors-resolver'
import components from './overrides'

export const theme = createTheme({
    variantColorResolver,
    components,
    cursorType: 'pointer',
    fontFamily:
        'Montserrat, Vazirmatn, Apple Color Emoji, Noto Sans SC, Twemoji Country Flags, sans-serif',
    fontFamilyMonospace: 'Fira Mono, monospace',
    breakpoints: {
        xs: '30em',
        sm: '40em',
        md: '48em',
        lg: '64em',
        xl: '80em',
        '2xl': '96em',
        '3xl': '120em',
        '4xl': '160em'
    },

    scale: 1,
    fontSmoothing: true,
    focusRing: 'never',
    white: '#ffffff',
    black: '#24292f',
    colors: {
        'exodus-logo-exo': [
            '#fcfcf9',
            '#f9f9f3',
            '#f6f6ee',
            '#f3f4e9',
            '#f1f2e4',
            '#eeefdf',
            '#eceddb',
            '#d0d1c1',
            '#b3b4a6',
            '#97988c'
        ],
        'exodus-logo-dus': [
            '#e8ecee',
            '#d7dde2',
            '#c3ccd3',
            '#b1bdc6',
            '#a0b0ba',
            '#90a2ae',
            '#8195a3',
            '#72838f',
            '#62717c',
            '#535f68'
        ],

        dark: [
            '#c9d1d9',
            '#b1bac4',
            '#8b949e',
            '#6e7681',
            '#484f58',
            '#30363d',
            '#21262d',
            '#161b22',
            '#0d1117',
            '#010409'
        ],
        'shaded-gray': [
            '#f5f5f5',
            '#e8e8e8',
            '#d4d4d4',
            '#c0c0c0',
            '#a8a8a8',
            '#a0a0a0',
            '#808080',
            '#686868',
            '#505050',
            '#383838'
        ]
    },
    primaryShade: 8,
    primaryColor: 'cyan',
    autoContrast: true,
    luminanceThreshold: 0.3,
    headings: {
        fontWeight: '600'
    },
    defaultRadius: 'md'
})
