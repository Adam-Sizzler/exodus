import { Monaco } from '@monaco-editor/react'
import {
    GetHostsCommand,
    SUBSCRIPTION_TEMPLATE_TYPE,
    TSubscriptionTemplateType
} from '@exodus/backend-contract'
import axios from 'axios'
import consola from 'consola'
import { configureMonacoYaml, MonacoYaml, MonacoYamlOptions, SchemasSettings } from 'monaco-yaml'
import { app } from 'src/config'

import { registerJsonSchema } from '@shared/utils/monaco/json-schema-registry'

type Host = GetHostsCommand.Response['response'][number]

interface IJsonSchemaDocument {
    [key: string]: unknown
    properties?: Record<string, unknown>
}

export const getTemplateModelPath = (templateType: TSubscriptionTemplateType) =>
    `subscription-template://${templateType.toLowerCase()}`

const YAML_OPTIONS: MonacoYamlOptions = {
    validate: true,
    enableSchemaRequest: true,
    hover: true,
    completion: true,
    format: {
        printWidth: 160,
        enable: true
    }
}

let monacoYaml: MonacoYaml | undefined

const configureYaml = (monaco: Monaco, schemas?: SchemasSettings[]) => {
    const options: MonacoYamlOptions = { ...YAML_OPTIONS, schemas }

    if (monacoYaml) {
        monacoYaml.update(options)
        return
    }

    monacoYaml = configureMonacoYaml(monaco, options)
}

const DOCS_URL = 'https://docs.ex/docs/learn/xray-json-advanced'
const DOCS_LINK = `\n\n[📖 Documentation](${DOCS_URL})`

const getInboundType = (host: Host): string => {
    return (host.inbound as any)?.type || 'unknown'
}

const formatHostTable = (hosts: Host[]): string => {
    const rows = [
        '| Host Name | Address | Port | Inbound Type |',
        '| :--- | :--- | :--- | :--- |'
    ]

    hosts.forEach((host) => {
        const name = host.remark || 'N/A'
        const address = host.address || 'N/A'
        const port = host.port || 'N/A'
        const inboundType = getInboundType(host)

        rows.push(`| **${name}** | \`${address}\` | \`${port}\` | \`${inboundType}\` |`)
    })

    return rows.join('\n')
}

export const configureMonaco = async (
    monaco: Monaco,
    language: 'json' | 'yaml',
    hosts: GetHostsCommand.Response['response'],
    templateType: TSubscriptionTemplateType
) => {
    try {
        if (language === 'yaml') {
            const schemas =
                templateType === SUBSCRIPTION_TEMPLATE_TYPE.MIHOMO
                    ? [
                          {
                              fileMatch: [getTemplateModelPath(templateType)],
                              uri: new URL(
                                  app.templateEditor.mihomoYamlSchemaUrl,
                                  window.location.origin
                              ).href
                          }
                      ]
                    : undefined

            configureYaml(monaco, schemas)
        }

        if (language === 'json') {
            const hostNames = hosts.map((h) => h.remark)
            const hostTable = formatHostTable(hosts)

            const schema = {
                type: 'object',
                properties: {
                    exodus: {
                        type: 'object',
                        markdownDescription: `Exodus-specific configuration options.${DOCS_LINK}`,
                        properties: {
                            hosts: {
                                type: 'array',
                                markdownDescription: `List of hosts to be included in the configuration.${DOCS_LINK}`,
                                items: {
                                    type: 'object',
                                    properties: {
                                        name: {
                                            type: 'string',
                                            enum: hostNames,
                                            markdownDescription: `### Available Hosts\n\n${hostTable}${DOCS_LINK}`
                                        },
                                        tag: {
                                            type: 'string',
                                            markdownDescription: `Custom tag for this host in the configuration.${DOCS_LINK}`
                                        },
                                        isExcluded: {
                                            type: 'boolean',
                                            markdownDescription: `Whether this host should be excluded from the configuration.${DOCS_LINK}`
                                        }
                                    }
                                }
                            },
                            outbounds: {
                                type: 'array',
                                markdownDescription: `Custom outbounds configuration.${DOCS_LINK}`,
                                items: {
                                    type: 'object',
                                    properties: {
                                        position: {
                                            type: 'string',
                                            enum: ['first', 'last'],
                                            markdownDescription: `Where to insert these outbounds in the configuration.${DOCS_LINK}`
                                        },
                                        target: {
                                            type: 'object',
                                            markdownDescription: `Target outbound configuration.${DOCS_LINK}`,
                                            properties: {
                                                tag: {
                                                    type: 'string',
                                                    markdownDescription: `Target outbound tag.${DOCS_LINK}`
                                                }
                                            }
                                        },
                                        outbounds: {
                                            type: 'array',
                                            markdownDescription: `List of outbounds to insert.${DOCS_LINK}`
                                        }
                                    }
                                }
                            },
                            rules: {
                                type: 'array',
                                markdownDescription: `Routing rules configuration.${DOCS_LINK}`,
                                items: {
                                    type: 'object',
                                    properties: {
                                        position: {
                                            type: 'string',
                                            enum: ['first', 'last'],
                                            markdownDescription: `Where to insert these rules in the routing configuration.${DOCS_LINK}`
                                        },
                                        target: {
                                            type: 'object',
                                            markdownDescription: `Target rule configuration.${DOCS_LINK}`,
                                            properties: {
                                                tag: {
                                                    type: 'string',
                                                    markdownDescription: `Target rule tag.${DOCS_LINK}`
                                                }
                                            }
                                        },
                                        rules: {
                                            type: 'array',
                                            markdownDescription: `List of rules to insert.${DOCS_LINK}`
                                        }
                                    }
                                }
                            },
                            balancers: {
                                type: 'array',
                                markdownDescription: `Load balancers configuration.${DOCS_LINK}`,
                                items: {
                                    type: 'object',
                                    properties: {
                                        position: {
                                            type: 'string',
                                            enum: ['first', 'last'],
                                            markdownDescription: `Where to insert these balancers in the routing configuration.${DOCS_LINK}`
                                        },
                                        target: {
                                            type: 'object',
                                            markdownDescription: `Target balancer configuration.${DOCS_LINK}`,
                                            properties: {
                                                tag: {
                                                    type: 'string',
                                                    markdownDescription: `Target balancer tag.${DOCS_LINK}`
                                                }
                                            }
                                        },
                                        balancers: {
                                            type: 'array',
                                            markdownDescription: `List of balancers to insert.${DOCS_LINK}`
                                        }
                                    }
                                }
                            },
                            reverse: {
                                type: 'object',
                                markdownDescription: `Reverse proxy configuration.${DOCS_LINK}`,
                                properties: {
                                    position: {
                                        type: 'string',
                                        enum: ['first', 'last'],
                                        markdownDescription: `Where to insert reverse proxy configuration.${DOCS_LINK}`
                                    },
                                    bridges: {
                                        type: 'array',
                                        markdownDescription: `Bridge configurations for reverse proxy.${DOCS_LINK}`
                                    },
                                    portals: {
                                        type: 'array',
                                        markdownDescription: `Portal configurations for reverse proxy.${DOCS_LINK}`
                                    }
                                }
                            },
                            fakeDns: {
                                type: 'object',
                                markdownDescription: `FakeDNS configuration.${DOCS_LINK}`,
                                properties: {
                                    position: {
                                        type: 'string',
                                        enum: ['first', 'last'],
                                        markdownDescription: `Where to insert FakeDNS configuration.${DOCS_LINK}`
                                    },
                                    pools: {
                                        type: 'array',
                                        markdownDescription: `List of FakeDNS pools.${DOCS_LINK}`,
                                        items: {
                                            type: 'object',
                                            properties: {
                                                ipPool: {
                                                    type: 'string',
                                                    markdownDescription: `IP pool CIDR for FakeDNS.${DOCS_LINK}`
                                                },
                                                poolSize: {
                                                    type: 'number',
                                                    markdownDescription: `Size of the FakeDNS pool.${DOCS_LINK}`
                                                }
                                            }
                                        }
                                    }
                                }
                            },
                            browserForwarder: {
                                type: 'object',
                                markdownDescription: `Browser forwarder configuration.${DOCS_LINK}`
                            },
                            observatory: {
                                type: 'object',
                                markdownDescription: `Observatory configuration.${DOCS_LINK}`
                            },
                            burstObservatory: {
                                type: 'object',
                                markdownDescription: `Burst observatory configuration.${DOCS_LINK}`
                            }
                        }
                    }
                }
            }

            registerJsonSchema({
                fileMatch: ['subscription-template://*', getTemplateModelPath(templateType)],
                schema,
                uri: 'https://subscription-template-schema.json'
            })

            if (templateType === SUBSCRIPTION_TEMPLATE_TYPE.SINGBOX) {
                const response = await axios.get<IJsonSchemaDocument>(
                    app.templateEditor.singboxJsonSchemaUrl
                )
                const singboxSchema = response.data

                registerJsonSchema({
                    fileMatch: [getTemplateModelPath(templateType)],
                    schema: {
                        ...singboxSchema,
                        properties: { ...singboxSchema.properties, exodus: true }
                    },
                    uri: 'https://singbox-schema.json'
                })
            }
        }
    } catch (error) {
        consola.error(`Failed to configure Monaco ${language.toUpperCase()}:`, error)
    }
}
