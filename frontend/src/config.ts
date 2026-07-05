import { withBasePath } from '@shared/constants/base-path'

export const app = {
    name: 'Exodus',
    githubRepo: 'https://github.com/Adam-Sizzler/backend',
    githubStars: 'https://github.com/Adam-Sizzler/backend/stargazers',
    githubIssues: 'https://github.com/Adam-Sizzler/backend/issues',
    githubOrg: 'https://github.com/Adam-Sizzler',
    githubDonation: 'https://github.com/Adam-Sizzler/backend#donation',
    configEditor: {
        wasmUrl: withBasePath('/assets/main.wasm'),
        wasmJsUrl: withBasePath('/assets/wasm_exec.js'),
        jsonSchemaUrl: withBasePath('/assets/singbox.schema.json'),
        jsonSchemaCnUrl: withBasePath('/assets/singbox.schema.json')
    }
}
