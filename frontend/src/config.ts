import { withBasePath } from '@shared/constants/base-path'

export const app = {
    name: 'Remnawave Dashboard',
    githubRepo: 'https://github.com/remnawave/backend',
    githubStars: 'https://github.com/remnawave/backend/stargazers',
    githubIssues: 'https://github.com/remnawave/backend/issues',
    githubOrg: 'https://github.com/remnawave',
    githubDonation: 'https://github.com/remnawave/backend#donation',
    configEditor: {
        wasmUrl: withBasePath('/assets/main.wasm'),
        wasmJsUrl: withBasePath('/assets/wasm_exec.js'),
        jsonSchemaUrl: withBasePath('/assets/singbox.schema.json'),
        jsonSchemaCnUrl: withBasePath('/assets/singbox.schema.json')
    }
}
