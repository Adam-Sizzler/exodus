/// <reference types="react" />
import type { Announcements, ScreenReaderInstructions } from '@dnd-kit/core/dist/components/Accessibility/types';
interface Props {
    announcements?: Announcements;
    container?: Element;
    screenReaderInstructions?: ScreenReaderInstructions;
    hiddenTextDescribedById: string;
}
export declare function Accessibility({ announcements, container, hiddenTextDescribedById, screenReaderInstructions, }: Props): JSX.Element | null;
export {};
