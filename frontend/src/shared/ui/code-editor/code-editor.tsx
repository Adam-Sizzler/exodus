import type { editor } from 'monaco-editor'

import { Box, Text } from '@mantine/core'
import { notifications } from '@mantine/notifications'
import Editor, { EditorProps, OnMount } from '@monaco-editor/react'
import clsx from 'clsx'
import { parse } from 'jsonc-parser'
import { Fragment, ReactNode, useRef, useState } from 'react'

import { BASE_MONACO_OPTIONS, MONACO_THEME_NAME } from '@shared/constants/monaco-theme'
import { describeJsonPath } from '@shared/utils/monaco/json-path'
import { RepairResult, repairJsonInEditor } from '@shared/utils/monaco/repair-json'

import { LoaderModalShared } from '../loader-modal/loader-model.shared'
import styles from './CodeEditor.module.css'

const REPAIR_ACTION_ID = 'exodus.repairJson'

const REPAIR_LABEL = 'Repair JSON'

const RESULT_COLORS: Record<RepairResult, { color: string; message: string }> = {
    failed: { color: 'red', message: 'Could not repair this JSON' },
    repaired: { color: 'teal', message: 'JSON repaired' },
    unchanged: { color: 'gray', message: 'Nothing to repair' }
}

interface Props extends Omit<EditorProps, 'wrapperProps'> {
    footer?: ReactNode
    withJsonPath?: boolean
    wrapperProps?: Record<string, unknown> & { className?: string }
}

export function CodeEditor(props: Props) {
    const {
        defaultLanguage,
        footer,
        language,
        onMount,
        options,
        withJsonPath,
        wrapperProps,
        ...rest
    } = props

    const [jsonPath, setJsonPath] = useState<string[]>([])
    const documentTextRef = useRef('')
    const documentValueRef = useRef<unknown>(undefined)

    const isJson = (language ?? defaultLanguage) === 'json'
    const showJsonPath = withJsonPath ?? isJson

    const syncDocumentSnapshot = (instance: editor.IStandaloneCodeEditor) => {
        const model = instance.getModel()

        if (!model) return

        const text = model.getValue()

        documentTextRef.current = text
        documentValueRef.current = parse(text)
    }

    const updateJsonPath = (instance: editor.IStandaloneCodeEditor) => {
        const model = instance.getModel()
        const [firstVisibleRange] = instance.getVisibleRanges()

        if (!model || !firstVisibleRange) {
            setJsonPath([])
            return
        }

        const topLine = firstVisibleRange.startLineNumber

        const offset = model.getOffsetAt({
            lineNumber: topLine,
            column: model.getLineMaxColumn(topLine)
        })

        const nextPath = describeJsonPath(documentTextRef.current, offset, documentValueRef.current)

        setJsonPath((currentPath) =>
            currentPath.length === nextPath.length &&
            currentPath.every((segment, index) => segment === nextPath[index])
                ? currentPath
                : nextPath
        )
    }

    const handleMount: OnMount = (instance, monaco) => {
        if (isJson) {
            instance.addAction({
                id: REPAIR_ACTION_ID,
                label: REPAIR_LABEL,
                contextMenuGroupId: '1_modification',
                contextMenuOrder: 1.32,
                run: (target) => {
                    const result = repairJsonInEditor(target as editor.IStandaloneCodeEditor)

                    notifications.show({
                        color: RESULT_COLORS[result].color,
                        message: RESULT_COLORS[result].message
                    })
                }
            })
        }

        syncDocumentSnapshot(instance)
        if (showJsonPath) updateJsonPath(instance)

        const changeDisposable = instance.onDidChangeModelContent(() => {
            syncDocumentSnapshot(instance)
            if (showJsonPath) updateJsonPath(instance)
        })

        const scrollDisposable = instance.onDidScrollChange(() => {
            if (showJsonPath) updateJsonPath(instance)
        })

        onMount?.(instance, monaco)

        return () => {
            changeDisposable.dispose()
            scrollDisposable.dispose()
        }
    }

    const mergedOptions: editor.IStandaloneEditorConstructionOptions = {
        ...BASE_MONACO_OPTIONS,
        ...options
    }

    return (
        <div className={styles.root}>
            {showJsonPath && (
                <div className={styles.pathBar}>
                    <Text c="dimmed" component="span" size="xs">
                        $
                    </Text>
                    {jsonPath.map((segment, index) => {
                        const isLast = index === jsonPath.length - 1

                        return (
                            <Fragment key={index}>
                                <Text
                                    className={styles.pathSeparator}
                                    component="span"
                                    mx={4}
                                    size="xs"
                                >
                                    .
                                </Text>
                                <Text
                                    className={clsx(isLast && styles.pathSegmentActive)}
                                    component="span"
                                    fw={isLast ? 600 : 400}
                                    size="xs"
                                >
                                    {segment}
                                </Text>
                            </Fragment>
                        )
                    })}
                </div>
            )}
            <Editor
                defaultLanguage={defaultLanguage}
                language={language}
                loading={<LoaderModalShared mih="100%" />}
                onMount={handleMount}
                options={mergedOptions}
                theme={MONACO_THEME_NAME}
                wrapperProps={{
                    ...wrapperProps,
                    className: clsx(styles.editorWrapper, wrapperProps?.className)
                }}
                {...rest}
            />

            {footer}
        </div>
    )
}

export { styles as editorClasses }
