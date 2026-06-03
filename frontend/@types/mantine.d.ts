import {
    ActionIconVariant,
    BadgeVariant,
    ButtonVariant,
    DefaultMantineColor,
    MantineColorsTuple,
    ThemeIconVariant
} from '@mantine/core'

type ExtendedThemeIconVariant = 'soft' | ThemeIconVariant
type ExtendedBadgeVariant = 'soft' | BadgeVariant

type ExtendedActionIconVariant = ExtendedThemeIconVariant | ActionIconVariant
type ExtendedButtonVariant = 'soft' | ButtonVariant
type ExtendedCustomColors = 'exodus' | 'shaded-gray' | DefaultMantineColor

declare module '@mantine/core' {
    export interface ThemeIconProps {
        variant?: ExtendedThemeIconVariant
    }

    export interface BadgeProps {
        variant?: ExtendedBadgeVariant
    }

    export interface ActionIconProps {
        variant?: ExtendedActionIconVariant
    }

    export interface ButtonProps {
        variant?: ExtendedButtonVariant
    }

    export interface MantineThemeColorsOverride {
        colors: Record<ExtendedCustomColors, MantineColorsTuple>
    }
}
