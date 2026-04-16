import { withBasePath } from '@shared/constants/base-path'

export const app = {
    name: 'Exodus Dashboard',
    githubRepo: 'https://github.com/teamdominant/exodus',
    githubStars: 'https://github.com/teamdominant/exodus/stargazers',
    githubIssues: 'https://github.com/teamdominant/exodus/issues',
    githubOrg: 'https://github.com/teamdominant',
    githubDonation: 'https://github.com/teamdominant/exodus',
    configEditor: {
        wasmUrl: withBasePath('/assets/main.wasm'),
        wasmJsUrl: withBasePath('/assets/wasm_exec.js'),
        jsonSchemaUrl: withBasePath('/assets/singbox.schema.json'),
        jsonSchemaCnUrl: withBasePath('/assets/singbox.schema.json')
    }
}
