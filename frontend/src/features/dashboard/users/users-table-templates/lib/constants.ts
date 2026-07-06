export const TEMPLATE_COLORS = [
    'violet',
    'blue',
    'teal',
    'green',
    'lime',
    'yellow',
    'orange',
    'red',
    'pink',
    'grape',
    'cyan'
] as const

export const DEFAULT_TEMPLATE_COLOR = TEMPLATE_COLORS[0]

export const getSwatchColor = (color: string) => `var(--mantine-color-${color}-6)`
