import { withBasePath } from '@shared/constants/base-path'

export const app = {
    name: 'Exodus',
    githubRepo: 'https://github.com/Adam-Sizzler/exodus',
    githubStars: 'https://github.com/Adam-Sizzler/exodus/stargazers',
    githubIssues: 'https://github.com/Adam-Sizzler/exodus/issues',
    githubTags: 'https://github.com/Adam-Sizzler/exodus/tags',
    githubOrg: 'https://github.com/Adam-Sizzler',
    githubDonation: 'https://github.com/Adam-Sizzler/exodus#donation',
    configEditor: {
        wasmUrl: withBasePath('/assets/main.wasm'),
        wasmJsUrl: withBasePath('/assets/wasm_exec.js'),
        jsonSchemaUrl: withBasePath('/assets/singbox.schema.json'),
        jsonSchemaCnUrl: withBasePath('/assets/singbox.schema.json')
    },
    templateEditor: {
        singboxJsonSchemaUrl: withBasePath('/assets/singbox.schema.json'),
        mihomoYamlSchemaUrl: withBasePath('/assets/mihomo.schema.json')
    }
}
