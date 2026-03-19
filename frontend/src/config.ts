import { withBasePath } from '@shared/constants/base-path'

export const app = {
    name: 'Cerberus Dashboard',
    githubRepo: 'https://github.com/cerberus/backend',
    githubStars: 'https://github.com/cerberus/backend/stargazers',
    githubIssues: 'https://github.com/cerberus/backend/issues',
    githubOrg: 'https://github.com/cerberus',
    githubDonation: 'https://github.com/cerberus/backend#donation',
    configEditor: {
        wasmUrl: withBasePath('/assets/main.wasm'),
        wasmJsUrl: withBasePath('/assets/wasm_exec.js'),
        jsonSchemaUrl: withBasePath('/assets/singbox.schema.json'),
        jsonSchemaCnUrl: withBasePath('/assets/singbox.schema.json')
    }
}
