import { ActionIcon, Box, Drawer, Flex, Group, Transition } from '@mantine/core'
import { useMediaQuery } from '@mantine/hooks'
import { GetConfigProfileByUuidCommand, GetSnippetsCommand } from '@exodus/backend-contract'
import { ConfigEditorWidget } from '@widgets/dashboard/config-profiles/config-editor/config-editor.widget'
import { SnippetsWidget } from '@widgets/dashboard/config-profiles/snippets-drawer/snippets.widget'
import { useTranslation } from 'react-i18next'
import { TbArrowBackUp, TbCode, TbFile } from 'react-icons/tb'
import { useNavigate } from 'react-router'

import { HelpActionIconShared } from '@shared/_modals/universal'
import { ROUTES } from '@shared/constants'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { Page } from '@shared/ui/page'
import { PageHeaderShared } from '@shared/ui/page-header/page-header.shared'

import { MODALS, useModalClose, useModalIsOpen } from '@entities/dashboard/modal-store'

interface Props {
    configProfile: GetConfigProfileByUuidCommand.Response['response']
    isWasmCrashed: boolean
    isWasmRestarting: boolean
    onRestartWasm: () => void
    snippets: GetSnippetsCommand.Response['response']
}

export const ConfigProfileByUuidPageComponent = (props: Props) => {
    const { configProfile, isWasmCrashed, isWasmRestarting, onRestartWasm, snippets } = props

    const { t } = useTranslation()
    const isMobile = useMediaQuery('(max-width: 1200px)')
    const navigate = useNavigate()

    const isOpen = useModalIsOpen(MODALS.CONFIG_PROFILE_SHOW_SNIPPETS_DRAWER)
    const close = useModalClose(MODALS.CONFIG_PROFILE_SHOW_SNIPPETS_DRAWER)

    return (
        <>
            <Page title={t('constants.config-profiles')}>
                <PageHeaderShared
                    actions={
                        <Group>
                            <HelpActionIconShared hidden={false} screen="PAGE_CONFIG_PROFILES" />

                            <ActionIcon
                                color="gray"
                                onClick={() =>
                                    navigate(ROUTES.DASHBOARD.MANAGEMENT.CONFIG_PROFILES)
                                }
                                size="input-md"
                                variant="soft"
                            >
                                <TbArrowBackUp size={24} />
                            </ActionIcon>
                        </Group>
                    }
                    description={configProfile.uuid}
                    icon={<TbFile size={24} />}
                    title={configProfile.name}
                />

                {isMobile ? (
                    <>
                        <ConfigEditorWidget
                            configProfile={configProfile}
                            isWasmCrashed={isWasmCrashed}
                            isWasmRestarting={isWasmRestarting}
                            onRestartWasm={onRestartWasm}
                            snippets={snippets}
                        />

                        <Drawer
                            keepMounted={false}
                            onClose={close}
                            opened={isOpen}
                            position="right"
                            size="450px"
                            title={
                                <BaseOverlayHeader
                                    iconColor="violet"
                                    IconComponent={TbCode}
                                    iconVariant="soft"
                                    title={t('snippets.drawer.widget.snippets')}
                                />
                            }
                            withCloseButton={true}
                        >
                            <SnippetsWidget />
                        </Drawer>
                    </>
                ) : (
                    <Flex gap="md">
                        <Box style={{ flex: 1, minWidth: 0 }}>
                            <ConfigEditorWidget
                                configProfile={configProfile}
                                isWasmCrashed={isWasmCrashed}
                                isWasmRestarting={isWasmRestarting}
                                onRestartWasm={onRestartWasm}
                                snippets={snippets}
                            />
                        </Box>

                        <Box
                            style={{
                                width: isOpen ? '400px' : '0px',
                                overflow: 'hidden',
                                transition: 'width 0.3s ease'
                            }}
                        >
                            <Transition
                                duration={300}
                                keepMounted
                                mounted={isOpen}
                                timingFunction="ease"
                                transition={{
                                    in: { opacity: 1, transform: 'translateX(0)' },
                                    out: { opacity: 0, transform: 'translateX(20px)' },
                                    transitionProperty: 'transform, opacity'
                                }}
                            >
                                {(transitionStyles) => (
                                    <Box
                                        style={{
                                            ...transitionStyles,
                                            width: '400px',
                                            pointerEvents: isOpen ? 'auto' : 'none'
                                        }}
                                    >
                                        <SnippetsWidget />
                                    </Box>
                                )}
                            </Transition>
                        </Box>
                    </Flex>
                )}
            </Page>
        </>
    )
}
