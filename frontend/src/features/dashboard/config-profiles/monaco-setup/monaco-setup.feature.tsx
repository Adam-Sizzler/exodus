import {
    GetSnippetsCommand,
    ResponseRulesConfigSchema,
    TSubscriptionTemplateType
} from '@exodus/backend-contract'
import zodToJsonSchema, { jsonDescription } from 'zod-to-json-schema'
import { Monaco } from '@monaco-editor/react'
import consola from 'consola'
import axios from 'axios'
import { z } from 'zod'

import { monacoTheme } from '@shared/constants/monaco-theme'
import { app } from 'src/config'

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

            const addSnippetProperty = (target?: { properties?: Record<string, unknown> }) => {
                if (!target) return
                target.properties = {
                    ...(target.properties ?? {}),
                    snippet: snippetSchema
                }
            }

            if (schema.definitions?.OutboundObject?.properties) {
                schema.definitions.OutboundObject.properties.snippet = snippetSchema
            }

            if (schema.definitions?.RuleObject?.properties) {
                schema.definitions.RuleObject.properties.snippet = snippetSchema
            }

            if (schema.definitions?.BalancerObject?.properties) {
                schema.definitions.BalancerObject.properties.snippet = snippetSchema
            }

            addSnippetProperty(schema.properties?.outbounds?.items)

            if (schema.properties?.route) {
                schema.properties.route.properties = {
                    ...(schema.properties.route.properties ?? {}),
                    rules: schema.properties.route.properties?.rules ?? {
                        type: 'array',
                        items: {
                            type: 'object',
                            additionalProperties: true
                        }
                    }
                }
                addSnippetProperty(schema.properties.route.properties.rules.items)
            }

            if (schema.properties?.routing) {
                schema.properties.routing.properties = {
                    ...(schema.properties.routing.properties ?? {}),
                    rules: schema.properties.routing.properties?.rules ?? {
                        type: 'array',
                        items: {
                            type: 'object',
                            additionalProperties: true
                        }
                    },
                    balancers: schema.properties.routing.properties?.balancers ?? {
                        type: 'array',
                        items: {
                            type: 'object',
                            additionalProperties: true
                        }
                    }
                }
                addSnippetProperty(schema.properties.routing.properties.rules.items)
                addSnippetProperty(schema.properties.routing.properties.balancers.items)
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
            const outboundObject =
                schema.definitions?.OutboundObject || schema.properties?.outbounds?.items || {
                    type: 'object',
                    additionalProperties: true
                }
            const ruleObject =
                schema.definitions?.RuleObject ||
                schema.properties?.route?.properties?.rules?.items || {
                    type: 'object',
                    additionalProperties: true
                }
            const balancerObject = schema.definitions?.BalancerObject || {
                type: 'object',
                additionalProperties: true
            }

            const snippetArraySchema = {
                $schema: 'http://json-schema.org/draft-07/schema#',
                title: 'Snippet Array',
                description: 'Array of Outbound, Rule or Balancer objects for snippets',
                type: 'array',
                items: {
                    oneOf: [
                        {
                            ...outboundObject,
                            title: 'Outbound Object',
                            description: 'Outbound configuration (for outbounds[])'
                        },
                        {
                            ...ruleObject,
                            title: 'Rule Object',
                            description: 'Routing rule (for routing.rules[])'
                        },
                        {
                            ...balancerObject,
                            title: 'Balancer Object',
                            description: 'Balancer configuration (for routing.balancers[])'
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

const ipListSchema = z.array(z.string()).default([])

export const NodePluginSchema = z
    .object({
        ingressFilter: z
            .object({
                enabled: z.boolean().default(false),
                blockedIps: ipListSchema
            })
            .optional()
            .describe('Ingress traffic filter.'),
        egressFilter: z
            .object({
                enabled: z.boolean().default(false),
                blockedIps: ipListSchema,
                blockedPorts: z.array(z.number()).default([])
            })
            .optional()
            .describe('Egress traffic filter.'),
        haproxyAuth: z
            .object({
                enabled: z.boolean().default(false)
            })
            .optional()
            .describe('Enable HAProxy authorization for the selected node.'),
        sharedLists: z
            .array(
                z.object({
                    name: z.string(),
                    type: z.literal('ipList'),
                    items: z.array(z.string())
                })
            )
            .default([])
    })
    .strict()

export const MonacoSetupNodePluginEditorFeature = {
    setup: async (monaco: Monaco) => {
        try {
            const schema = zodToJsonSchema(NodePluginSchema, {
                name: 'Node Plugin Config Schema',
                applyRegexFlags: true,
                errorMessages: true,
                postProcess: jsonDescription
            })

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
            consola.error('Failed to load node plugin JSON schema:', error)
        }
    }
}
