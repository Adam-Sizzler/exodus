import type { editor } from 'monaco-editor'

import {
    TbClipboardCopy,
    TbClipboardText,
    TbCut,
    TbDownload,
    TbMenuDeep,
    TbSelectAll
} from 'react-icons/tb'
import { useClipboard, useDisclosure, useMediaQuery } from '@mantine/hooks'
import { PiCheckSquareOffset, PiFloppyDisk } from 'react-icons/pi'
import { ActionIcon, Button, Group, Menu } from '@mantine/core'
import { notifications } from '@mantine/notifications'
import { Monaco } from '@monaco-editor/react'
import consola from 'consola/browser'
import { RefObject } from 'react'

import { useDownloadTemplate } from '@shared/ui/load-templates/use-download-template'
import { QueryKeys, useUpdateNodePlugin } from '@shared/api/hooks'
import { queryClient } from '@shared/api'

interface Props {
    editorRef: RefObject<editor.IStandaloneCodeEditor | null>
    hasUnsavedChanges: boolean
    isNodePluginValid: boolean
    monacoRef: RefObject<Monaco | null>
    originalValue: string
    pluginUuid: string
    setHasUnsavedChanges: (value: boolean) => void
    setIsNodePluginValid: (value: boolean) => void
    setOriginalValue: (value: string) => void
    setResult: (value: string) => void
}

const unsupportedPluginKeys = ['torrentBlocker', 'connectionDrop']

export function NodePluginsEditorActionsFeature(props: Props) {
    const {
        editorRef,
        monacoRef,
        setResult,
        setIsNodePluginValid,
        hasUnsavedChanges,
        setHasUnsavedChanges,
        setOriginalValue,
        pluginUuid
    } = props

    const isMobile = useMediaQuery('(max-width: 48em)')
    const clipboard = useClipboard({ timeout: 500 })
    const [opened, handlers] = useDisclosure(false)

    const { mutate: updateNodePlugin, isPending: isUpdating } = useUpdateNodePlugin({
        mutationFns: {
            onSuccess: async (updatedNodePlugin) => {
                await queryClient.refetchQueries({
                    queryKey: QueryKeys.nodes.getNodePlugin({ uuid: pluginUuid }).queryKey
                })

                setIsNodePluginValid(true)

                const newValue = JSON.stringify(updatedNodePlugin.pluginConfig, null, 2)

                if (editorRef.current) {
                    editorRef.current.setValue(newValue)
                    setOriginalValue(newValue)
                }

                await queryClient.setQueryData(
                    QueryKeys.nodes.getNodePlugin({ uuid: pluginUuid }).queryKey,
                    updatedNodePlugin
                )

                setHasUnsavedChanges(false)
            },
            onError: (error) => {
                setIsNodePluginValid(false)
                setResult(error.message)
            }
        }
    })

    const { openDownloadModal } = useDownloadTemplate({
        editorType: 'NODE_PLUGIN',
        templateType: 'NODE_PLUGIN',
        editorRef
    })

    const handleSave = () => {
        if (!editorRef.current) return
        if (!monacoRef.current) return

        const currentValue = editorRef.current.getValue()

        let parsed: Record<string, unknown>
        try {
            parsed = JSON.parse(currentValue) as Record<string, unknown>
        } catch (error) {
            consola.error(error)
            notifications.show({
                color: 'red',
                message: 'Failed to save invalid JSON.',
                title: 'Error'
            })
            return
        }

        const unsupportedKey = unsupportedPluginKeys.find((key) => key in parsed)
        if (unsupportedKey) {
            notifications.show({
                color: 'red',
                message: `${unsupportedKey} is not supported by sing-box core.`,
                title: 'Error'
            })
            return
        }

        updateNodePlugin({
            route: {
                uuid: pluginUuid
            },
            variables: {
                pluginConfig: parsed
            }
        })
    }

    const handleCopyConfig = () => {
        if (!editorRef.current) return
        clipboard.copy(editorRef.current.getValue())
    }

    const handleSelectAll = () => {
        if (!editorRef.current) return

        const model = editorRef.current.getModel()
        if (!model) return

        editorRef.current.setSelection({
            startLineNumber: 1,
            startColumn: 1,
            endLineNumber: model.getLineCount(),
            endColumn: model.getLineMaxColumn(model.getLineCount())
        })
    }

    const handleCut = () => {
        if (!editorRef.current) return

        const selection = editorRef.current.getSelection()
        const model = editorRef.current.getModel()
        if (!selection || !model) return

        const selectedText = model.getValueInRange(selection)
        clipboard.copy(selectedText)

        editorRef.current.executeEdits('', [{ range: selection, text: '' }])
    }

    const handlePaste = () => {
        if (!editorRef.current) return

        const position = editorRef.current.getPosition()
        if (!position) return

        navigator.clipboard.readText().then((text) => {
            if (!editorRef.current) return
            editorRef.current.executeEdits('', [
                {
                    range: {
                        startLineNumber: position.lineNumber,
                        startColumn: position.column,
                        endLineNumber: position.lineNumber,
                        endColumn: position.column
                    },
                    text
                }
            ])
        })
    }

    const formatDocument = () => {
        if (!editorRef.current) return
        editorRef.current.getAction('editor.action.formatDocument')?.run()
    }

    return (
        <Group grow={isMobile} preventGrowOverflow={false} wrap="wrap">
            <Button
                color={!hasUnsavedChanges ? 'gray' : 'teal'}
                disabled={!hasUnsavedChanges}
                leftSection={<PiFloppyDisk size={16} />}
                loading={isUpdating}
                onClick={handleSave}
            >
                Save
            </Button>

            <Group gap={0} wrap="nowrap">
                <Menu
                    onClose={() => handlers.close()}
                    onOpen={() => handlers.open()}
                    shadow="md"
                    trigger="click-hover"
                    withinPortal
                >
                    <Menu.Target>
                        <ActionIcon
                            size={36}
                            style={{
                                borderTopRightRadius: 0,
                                borderBottomRightRadius: 0
                            }}
                            variant={opened ? 'outline' : 'default'}
                        >
                            <TbMenuDeep size={20} />
                        </ActionIcon>
                    </Menu.Target>

                    <Menu.Dropdown>
                        <Menu.Item
                            color={clipboard.copied ? 'teal' : undefined}
                            leftSection={<TbClipboardCopy size={14} />}
                            onClick={handleCopyConfig}
                        >
                            Copy all content
                        </Menu.Item>

                        <Menu.Item
                            leftSection={<TbSelectAll size={14} />}
                            onClick={handleSelectAll}
                        >
                            Select all
                        </Menu.Item>

                        <Menu.Item leftSection={<TbCut size={14} />} onClick={handleCut}>
                            Cut selection
                        </Menu.Item>

                        <Menu.Item
                            leftSection={<TbClipboardText size={14} />}
                            onClick={handlePaste}
                        >
                            Paste from clipboard
                        </Menu.Item>

                        <Menu.Divider />

                        <Menu.Item leftSection={<TbDownload size={14} />} onClick={openDownloadModal}>
                            Load from GitHub
                        </Menu.Item>
                    </Menu.Dropdown>
                </Menu>

                <Button
                    leftSection={<PiCheckSquareOffset size={16} />}
                    onClick={formatDocument}
                    style={{
                        borderTopLeftRadius: 0,
                        borderBottomLeftRadius: 0,
                        borderLeft: 0,
                        width: '100%'
                    }}
                    variant="default"
                >
                    Format
                </Button>
            </Group>
        </Group>
    )
}
