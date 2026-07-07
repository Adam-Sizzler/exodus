import { withBasePath } from '@shared/constants/base-path'
export const isPwa = () => window.matchMedia('(display-mode: standalone)').matches

export const openOrNavigate = (url: string, navigate: (url: string) => void) => {
    if (isPwa()) {
        navigate(url)
    } else {
        window.open(withBasePath(url), '_blank')
    }
}
