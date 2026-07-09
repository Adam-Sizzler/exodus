import {
    GetSnippetsCommand,
    ResponseRulesConfigSchema,
    TSubscriptionTemplateType
} from '@exodus/backend-contract'
import { NodePluginSchema } from '@exodus/node-plugins'
import { Monaco } from '@monaco-editor/react'
import axios from 'axios'
import consola from 'consola'
import { app } from 'src/config'
import zodToJsonSchema, { jsonDescription } from 'zod-to-json-schema'

import { monacoTheme } from '@shared/constants/monaco-theme'
import type { NodePluginHaproxyInboundTagOption } from '@widgets/dashboard/node-plugins/node-plugin-editor/node-plugin-editor-schema'
import { HAPROXY_AUTH_ALL_INBOUNDS_TAG } from '@widgets/dashboard/node-plugins/node-plugin-editor/node-plugin-editor-schema'

export const MonacoSetupFeature = {
    setup: async (
        monaco: Monaco,
        currentLanguage: string,
        snippets: GetSnippetsCommand.Response['response']['snippets']
    ) => {
        try {
            const snippetNames = snippets.map((s) => s.name)

            let { jsonSchemaUrl } = app.configEditor
            switch (currentLanguage) {
                case 'zh':
                    jsonSchemaUrl = app.configEditor.jsonSchemaCnUrl
                    break
                default:
                    jsonSchemaUrl = app.configEditor.jsonSchemaUrl
            }

            const response = await axios.get(jsonSchemaUrl)
            const schema = await response.data

            const snippetDescriptions = snippets.map((snippet) => {
                const snippetJson = JSON.stringify(snippet.snippet, null, 1)

                return ['', '```json', snippetJson.slice(2, -2), '```', '', '---', ''].join('\n')
            })

            const snippetSchema = {
                name: 'snippet',
                title: 'Exodus Snippets',
                markdownDescription:
                    'Create your own snippets to quickly configure your **Outbounds** or **Rules**. \n\n\nReference them here, Exodus will handle the rest.',
                type: 'string',
                enum: snippetNames,
                markdownEnumDescriptions: snippetDescriptions,
                minLength: 2,
                maxLength: 255,
                pattern: '^[A-Za-z0-9_\\s-]+$',
                patternErrorMessage:
                    'Snippet name can only contain: letters, numbers, spaces, _ and -'
            }

            if (schema.definitions?.OutboundObject?.properties) {
                schema.definitions.OutboundObject.properties.snippet = snippetSchema
            }

            if (schema.definitions?.RuleObject?.properties) {
                schema.definitions.RuleObject.properties.snippet = snippetSchema
            }

            if (schema.properties?.outbounds?.items) {
                schema.properties.outbounds.items.properties = {
                    ...schema.properties.outbounds.items.properties,
                    snippet: snippetSchema
                }
            }

            if (schema.properties?.route) {
                schema.properties.route.properties = schema.properties.route.properties ?? {}
                const routeProperties = schema.properties.route.properties

                routeProperties.rules = routeProperties.rules ?? {
                    type: 'array',
                    items: {
                        type: 'object',
                        additionalProperties: true
                    }
                }

                routeProperties.rules.items = routeProperties.rules.items ?? {
                    type: 'object',
                    additionalProperties: true
                }
                routeProperties.rules.items.properties = {
                    ...routeProperties.rules.items.properties,
                    snippet: snippetSchema
                }
            }

            monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
                allowComments: false,
                enableSchemaRequest: true,
                schemaRequest: 'warning',
                schemas: [
                    {
                        fileMatch: ['*'],
                        schema,
                        uri: 'https://singbox-config-schema.json'
                    }
                ],
                validate: true
            })
        } catch (error) {
            consola.error('Failed to load JSON schema:', error)
        }
    }
}
export const MonacoSetupSnippetsFeature = {
    setup: async (monaco: Monaco, currentLanguage: string) => {
        try {
            let { jsonSchemaUrl } = app.configEditor
            switch (currentLanguage) {
                case 'zh':
                    jsonSchemaUrl = app.configEditor.jsonSchemaCnUrl
                    break
                default:
                    jsonSchemaUrl = app.configEditor.jsonSchemaUrl
            }

            const response = await axios.get(jsonSchemaUrl)
            const schema = await response.data

            const snippetArraySchema = {
                $schema: 'http://json-schema.org/draft-07/schema#',
                title: 'Snippet Array',
                description: 'Array of Outbound, Rule or Rule Set objects for snippets',
                type: 'array',
                items: {
                    oneOf: [
                        {
                            ...schema.definitions?.OutboundObject,
                            title: 'Outbound Object',
                            description: 'Outbound configuration (for outbounds[])'
                        },
                        {
                            ...schema.definitions?.RuleObject,
                            title: 'Rule Object',
                            description: 'Routing rule (for route.rules[])'
                        },
                        {
                            ...schema.definitions?.RuleSetObject,
                            title: 'Rule Set Object',
                            description: 'Rule set configuration (for route.rule_set[])'
                        }
                    ]
                },
                minItems: 1,
                definitions: schema.definitions || {}
            }

            monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
                allowComments: false,
                enableSchemaRequest: true,
                schemaRequest: 'warning',
                schemas: [
                    {
                        fileMatch: ['snippet://*'],
                        schema: snippetArraySchema,
                        uri: 'https://snippet-schema.json'
                    }
                ],
                validate: true
            })

            return snippetArraySchema
        } catch (error) {
            consola.error('Failed to load snippet JSON schema:', error)
            return null
        }
    }
}

export const MonacoSetupResponseRulesFeature = {
    setup: async (
        monaco: Monaco,
        groupedTemplates: Record<TSubscriptionTemplateType, string[]>
    ) => {
        try {
            const schema = zodToJsonSchema(ResponseRulesConfigSchema, {
                name: 'Response Rules Config Schema',
                applyRegexFlags: true,
                errorMessages: true,
                postProcess: jsonDescription
            })

            const templateOptions = {
                BROWSER: [],
                BLOCK: [],
                STATUS_CODE_404: [],
                STATUS_CODE_451: [],
                SOCKET_DROP: [],
                ...groupedTemplates
            }

            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const schemaDefinitions = (schema as any).definitions?.['Response Rules Config Schema']
            const rulesItems = schemaDefinitions?.properties?.rules?.items

            if (rulesItems) {
                if (!rulesItems.allOf) {
                    rulesItems.allOf = []
                }

                Object.entries(templateOptions).forEach(([responseType, templates]) => {
                    if (templates.length > 0) {
                        rulesItems.allOf.push({
                            if: {
                                properties: {
                                    responseType: { const: responseType }
                                },
                                required: ['responseType']
                            },
                            // oxlint-disable-next-line
                            then: {
                                properties: {
                                    responseModifications: {
                                        properties: {
                                            subscriptionTemplate: {
                                                enum: templates,
                                                markdownDescription: `Available templates for **${responseType}** response type.`,
                                                markdownEnumDescriptions: templates.map(
                                                    (t) => `Use ${t} template`
                                                )
                                            }
                                        }
                                    }
                                }
                            }
                        })
                    } else {
                        rulesItems.allOf.push({
                            if: {
                                properties: {
                                    responseType: { const: responseType }
                                },
                                required: ['responseType']
                            },
                            // oxlint-disable-next-line
                            then: {
                                properties: {
                                    responseModifications: {
                                        properties: {
                                            subscriptionTemplate: {
                                                type: 'null',
                                                not: { type: 'string' },
                                                markdownDescription: `⚠️ No templates available for **${responseType}** response type. This field should not be used.`
                                            }
                                        }
                                    }
                                }
                            }
                        })
                    }
                })
                //     } else {
                //         rulesItems.allOf.push({
                //             if: {
                //                 properties: {
                //                     responseType: { const: responseType }
                //                 },
                //                 required: ['responseType']
                //             },
                //             then: {
                //                 properties: {
                //                     responseModifications: {
                //                         properties: {
                //                             overrideSubscriptionTemplateWith: false
                //                         }
                //                     }
                //                 }
                //             }
                //         })
                //     }
                // })
            }

            monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
                schemaValidation: 'error',
                comments: 'error',
                trailingCommas: 'error',

                schemas: [
                    {
                        fileMatch: ['response-rules://*'],
                        schema,
                        uri: 'https://response-rules-schema.json'
                    }
                ],
                validate: true
            })

            monaco.languages.json.jsonDefaults.setModeConfiguration({
                documentFormattingEdits: true,
                documentRangeFormattingEdits: true,
                completionItems: true,
                hovers: true,
                documentSymbols: true,
                tokens: true,
                colors: true,
                foldingRanges: true,
                diagnostics: true,
                selectionRanges: true
            })

            monaco.editor.defineTheme('GithubDark', {
                ...monacoTheme,
                base: 'vs-dark'
            })
        } catch (error) {
            consola.error('Failed to load JSON schema:', error)
        }
    }
}

export const MonacoSetupNodePluginEditorFeature = {
    setup: async (
        monaco: Monaco,
        haproxyInboundTagOptions?: NodePluginHaproxyInboundTagOption[]
    ) => {
        try {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const schema = zodToJsonSchema(NodePluginSchema, {
                name: 'Node Plugin Schema',
                applyRegexFlags: true,
                errorMessages: true,
                postProcess: jsonDescription
            }) as any

            const nodePluginDefinition = schema.definitions?.['Node Plugin Schema']
            const inboundTagsSchema = nodePluginDefinition?.properties?.haproxyAuth?.properties?.inboundTags

            if (inboundTagsSchema && haproxyInboundTagOptions) {
                const tags = Array.from(
                    new Set(
                        haproxyInboundTagOptions
                            .map((option) => option.tag.trim())
                            .filter((tag) => tag.length > 0 && tag !== HAPROXY_AUTH_ALL_INBOUNDS_TAG)
                    )
                ).sort((a, b) => a.localeCompare(b))

                const descriptionsByTag = haproxyInboundTagOptions.reduce<Record<string, string>>(
                    (acc, option) => {
                        const tag = option.tag.trim()
                        if (!tag) return acc

                        const nodeNames = option.nodeNames.filter(Boolean).join(', ')
                        const metadata = [
                            option.type ? `type: ${option.type}` : null,
                            nodeNames ? `nodes: ${nodeNames}` : null
                        ]
                            .filter(Boolean)
                            .join('; ')

                        acc[tag] = metadata
                            ? `Inbound tag \`${tag}\` (${metadata}).`
                            : `Inbound tag \`${tag}\`.`
                        return acc
                    },
                    {}
                )

                inboundTagsSchema.items = {
                    ...inboundTagsSchema.items,
                    type: 'string',
                    enum: [HAPROXY_AUTH_ALL_INBOUNDS_TAG, ...tags],
                    markdownEnumDescriptions: [
                        'All supported inbounds assigned to the node.',
                        ...tags.map((tag) => descriptionsByTag[tag] ?? `Inbound tag \`${tag}\`.`)
                    ]
                }
            }
            monaco.languages.json.jsonDefaults.setDiagnosticsOptions({
                schemaValidation: 'error',
                comments: 'error',
                trailingCommas: 'error',

                schemas: [
                    {
                        fileMatch: ['node-plugin://*'],
                        schema,
                        uri: 'https://node-plugin-schema.json'
                    }
                ],
                validate: true
            })

            monaco.languages.json.jsonDefaults.setModeConfiguration({
                documentFormattingEdits: true,
                documentRangeFormattingEdits: true,
                completionItems: true,
                hovers: true,
                documentSymbols: true,
                tokens: true,
                colors: true,
                foldingRanges: true,
                diagnostics: true,
                selectionRanges: true
            })

            monaco.editor.defineTheme('GithubDark', {
                ...monacoTheme,
                base: 'vs-dark'
            })
        } catch (error) {
            consola.error('Failed to load JSON schema:', error)
        }
    }
}
