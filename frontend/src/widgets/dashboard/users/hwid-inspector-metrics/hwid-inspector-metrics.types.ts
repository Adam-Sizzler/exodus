export interface IPlatformApp {
    app: string
    count: number
}

export interface IPlatformDatum {
    allApps: IPlatformApp[]
    apps: IPlatformApp[]
    color: string
    name: string
    value: number
}
