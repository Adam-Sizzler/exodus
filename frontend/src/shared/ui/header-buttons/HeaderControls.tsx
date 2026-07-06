import { BoxProps, Group } from '@mantine/core'

import { GithubControl } from './GithubControl'
import { LanguageControl } from './LanguageControl'
import { LogoutControl } from './LogoutControl'
import { PrimeControl } from './PrimeControl'
import { RecapControl } from './RecapControl'
import { SupportControl } from './SupportControl'
import { TelegramControl } from './TelegramControl'
import { VersionControl } from './VersionControl'

interface HeaderControlsProps extends BoxProps {
    githubLink?: string
    isGithubLoading?: boolean
    stars?: number
    telegramLink: string
    withGithub?: boolean
    withLanguage?: boolean
    withLogout?: boolean
    withPrime?: boolean
    withRecap?: boolean
    withSupport?: boolean
    withTelegram?: boolean
    withVersion?: boolean
}

export function HeaderControls({
    githubLink,
    withGithub = true,
    withTelegram = true,
    withSupport = true,
    withLogout = true,
    withLanguage = true,
    withVersion = true,
    withRecap = false,
    withPrime = false,
    telegramLink,
    stars,
    isGithubLoading,
    ...others
}: HeaderControlsProps) {
    return (
        <Group gap="xs" {...others}>
            {withTelegram && <TelegramControl link={telegramLink} />}
            {withPrime && <PrimeControl />}
            {withSupport && <SupportControl />}

            {withVersion && <VersionControl />}
            {withGithub && (
                <GithubControl isLoading={isGithubLoading} link={githubLink!} stars={stars} />
            )}
            {withRecap && <RecapControl />}
            {withLanguage && <LanguageControl />}
            {withLogout && <LogoutControl />}
        </Group>
    )
}
