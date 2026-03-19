import { Box, Center, Code, Drawer, Stack, Title, Typography } from '@mantine/core'
import { TbAlertCircle, TbQuestionMark } from 'react-icons/tb'
import { memo, useEffect, useState } from 'react'
import { useTranslation } from 'node_modules/react-i18next'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeRaw from 'rehype-raw'

import { MODALS, useModalClose, useModalState } from '@entities/dashboard/modal-store'

import { BaseOverlayHeader } from '../overlays/base-overlay-header'
import { THelpDrawerAvailableScreen } from './help-drawer.types'
import { LoaderModalShared } from '../loader-modal'
import classes from './help-drawer.module.css'

const SUPPORTED_LANGUAGES = new Set(['en', 'fa', 'ru', 'zh'])

const resolveDocsUrl = (screen: THelpDrawerAvailableScreen, language: string) => {
    const lang = language.split('-')[0]
    const safeLang = SUPPORTED_LANGUAGES.has(lang) ? lang : 'en'
    return `https://raw.githubusercontent.com/cerberus/panel/refs/heads/main/_panel-docs/help-articles/${safeLang}/${screen}.md`
}

const getStaticHelpArticle = (screen: THelpDrawerAvailableScreen, language: string): null | string => {
    const lang = language.split('-')[0]

    if (screen === 'MODULES_HAPROXY') {
        if (lang === 'ru') {
            return `## Модуль: HAPROXY

Кратко: модуль формирует таблицу учетных данных для Lua-авторизации в HAProxy и поддерживает ее в актуальном состоянии при деплое ноды.

### Что делает модуль

- Создает и обновляет **один** файл пользователей.
- Добавляет в файл всех пользователей текущей ноды (без привязки к конкретному inbound).
- Для каждого пользователя пишет:
- \`VLESS UUID\`
- \`Trojan credential\` (SHA-224 hash, 56 символов).
- Поле состояния всегда \`1\` (включен).

### Какой файл создается

- Относительный путь в модуле: \`haproxy/data/users.csv\`
- Абсолютный путь в контейнере ноды: \`/app/haproxy/data/users.csv\`

### Формат строк

\`1,username,credential\`

### Пример содержимого

\`\`\`
1,pablo_escobar,90cdd126-819d-5f2f-a1dd-4cbf085182b9
1,alvarez,9f37a85cafb66ab9c36f866e5c787b2e42a997716c02337f532a259c
1,guest,8f49d237686bee189aaedaba1ad434c5580356d8e8b5a3d702bd0228
\`\`\`

### Для чего это нужно

- \`haproxy/lua/auth.lua\` читает пользователей из \`/etc/haproxy/data/users.csv\` и определяет протокол:
- Trojan (включая проверку MUX)
- VLESS (по UUID в бинарном рукопожатии)
- \`haproxy/lua/users_loader.lua\` разбирает CSV и кладет значения в карты:
- \`trojan[sha224] -> username\`
- \`vless[uuid] -> username\`
- Дальше \`haproxy.cfg\` использует \`lua.identify_protocol\` и направляет трафик в backend \`trojan\` или \`vless\`.

### Важно

- Модуль не изменяет \`haproxy.cfg\` и Lua-скрипты.
- Модуль управляет только файлом \`users.csv\`.`
        }

        return `## Module: HAPROXY

In short: this module builds and maintains credentials table used by HAProxy Lua auth during node deploy.

### What this module does

- Creates and updates **one** file with users.
- Includes all users of the current node (not tied to specific inbound).
- Writes for each user:
- \`VLESS UUID\`
- \`Trojan credential\` (SHA-224 hash, 56 chars).
- Enabled state is always \`1\`.

### File path

- Relative path in module logic: \`haproxy/data/users.csv\`
- Absolute path inside node container: \`/app/haproxy/data/users.csv\`

### Row format

\`1,username,credential\`

### Example

\`\`\`
1,pablo_escobar,90cdd126-819d-5f2f-a1dd-4cbf085182b9
1,alvarez,9f37a85cafb66ab9c36f866e5c787b2e42a997716c02337f532a259c
1,guest,8f49d237686bee189aaedaba1ad434c5580356d8e8b5a3d702bd0228
\`\`\`

### Why it is needed

- \`haproxy/lua/auth.lua\` loads users from \`/etc/haproxy/data/users.csv\` and detects protocol:
- Trojan (with MUX check)
- VLESS (UUID parsed from binary handshake)
- \`haproxy/lua/users_loader.lua\` parses CSV into maps:
- \`trojan[sha224] -> username\`
- \`vless[uuid] -> username\`
- Then \`haproxy.cfg\` uses \`lua.identify_protocol\` and routes traffic to \`trojan\` or \`vless\` backend.

### Important

- The module does not modify \`haproxy.cfg\` or Lua scripts.
- It manages only \`users.csv\`.`
    }

    return null
}

export const HelpDrawerShared = memo(() => {
    const { t, i18n } = useTranslation()

    const { isOpen, internalState } = useModalState(MODALS.HELP_DRAWER)
    const close = useModalClose(MODALS.HELP_DRAWER)

    const [content, setContent] = useState('')
    const [loading, setLoading] = useState(false)
    const [showContent, setShowContent] = useState(false)
    const [error, setError] = useState<null | string>(null)

    useEffect(() => {
        if (!isOpen || !internalState) {
            return
        }

        setLoading(true)
        setError(null)
        setContent('')
        setShowContent(false)

        const staticContent = getStaticHelpArticle(internalState.screen, i18n.language)
        if (staticContent) {
            setContent(staticContent)
            setTimeout(() => {
                setLoading(false)
                setShowContent(true)
            }, 150)
            return
        }

        fetch(resolveDocsUrl(internalState.screen, i18n.language))
            .then((res) => {
                if (!res.ok) throw new Error(t('help-drawer.shared.failed-to-load-documentation'))
                return res.text()
            })
            .then((text) => {
                setContent(text)
            })
            .catch((err) => {
                setError(err.message)
                setLoading(false)
            })
            .finally(() => {
                setTimeout(() => {
                    setLoading(false)
                    setShowContent(true)
                }, 300)
            })
    }, [isOpen])

    const cleanContent = () => {
        setContent('')
        setShowContent(false)
        setLoading(false)
        setError(null)
    }

    return (
        <Drawer
            keepMounted={false}
            onClose={close}
            onExitTransitionEnd={cleanContent}
            opened={isOpen}
            overlayProps={{ backgroundOpacity: 0.6, blur: 0 }}
            position="right"
            size="lg"
            title={
                <BaseOverlayHeader
                    IconComponent={TbQuestionMark}
                    iconVariant="gradient-yellow"
                    title={t('help-action-icon.shared.help-article')}
                />
            }
        >
            {loading && (
                <LoaderModalShared
                    h="80vh"
                    text={t('help-drawer.shared.loading-documentation')}
                    w="100%"
                />
            )}

            {error && (
                <Center h="80vh" w="100%">
                    <Stack align="center" gap="xs">
                        <TbAlertCircle color="var(--mantine-color-red-5)" size="4rem" />
                        <Title c="dimmed" order={4} size="lg">
                            {t('help-drawer.shared.failed-to-load-documentation')}
                        </Title>
                        <Code color="var(--mantine-color-red-light)">{error}</Code>
                    </Stack>
                </Center>
            )}

            <Box
                style={{
                    opacity: showContent ? 1 : 0,
                    pointerEvents: showContent ? 'auto' : 'none',
                    transition: 'opacity 0.3s ease'
                }}
            >
                {!loading && content && (
                    <Typography className={classes.root}>
                        <ReactMarkdown rehypePlugins={[rehypeRaw]} remarkPlugins={[remarkGfm]}>
                            {content}
                        </ReactMarkdown>
                    </Typography>
                )}
            </Box>
        </Drawer>
    )
})
