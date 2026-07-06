export interface IActions {
    actions: {
        setDisclaimerAccepted: (accepted: boolean) => void
        setMobileWarningClosed: (closed: boolean) => void
        setSrrAdvancedModalClosed: (closed: boolean) => void
    }
}
