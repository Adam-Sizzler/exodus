import { ActionIconVariant, BadgeVariant, ButtonVariant, ThemeIconVariant } from '@mantine/core'

type ExtendedThemeIconVariant =
    | 'soft'
    | 'gradient-blue'
    | 'gradient-cyan'
    | 'gradient-gray'
    | 'gradient-green'
    | 'gradient-indigo'
    | 'gradient-lime'
    | 'gradient-orange'
    | 'gradient-pink'
    | 'gradient-red'
    | 'gradient-teal'
    | 'gradient-violet'
    | 'gradient-yellow'
    | ThemeIconVariant

type ExtendedBadgeVariant =
    | 'soft'
    | 'gradient-blue'
    | 'gradient-cyan'
    | 'gradient-gray'
    | 'gradient-green'
    | 'gradient-indigo'
    | 'gradient-lime'
    | 'gradient-orange'
    | 'gradient-pink'
    | 'gradient-red'
    | 'gradient-teal'
    | 'gradient-violet'
    | 'gradient-yellow'
    | BadgeVariant

type ExtendedActionIconVariant = ExtendedThemeIconVariant | ActionIconVariant
type ExtendedButtonVariant = ExtendedThemeIconVariant | ButtonVariant

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
}
