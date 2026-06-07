import { withBasePath } from '@shared/constants/base-path'

export const app = {
    name: 'Exodus Dashboard',
    githubRepo: 'https://github.com/Adam-Sizzler/exodus',
    githubStars: 'https://github.com/Adam-Sizzler/exodus/stargazers',
    githubIssues: 'https://github.com/Adam-Sizzler/exodus/issues',
    githubOrg: 'https://github.com/Adam-Sizzler',
    githubDonation: 'https://github.com/Adam-Sizzler/exodus',
    configEditor: {
        wasmUrl: withBasePath('/assets/main.wasm'),
        wasmJsUrl: withBasePath('/assets/wasm_exec.js'),
        jsonSchemaUrl: withBasePath('/assets/singbox.schema.json'),
        jsonSchemaCnUrl: withBasePath('/assets/singbox.schema.json')
    }
}
