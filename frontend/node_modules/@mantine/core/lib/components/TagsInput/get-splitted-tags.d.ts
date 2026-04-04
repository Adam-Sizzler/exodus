interface GetSplittedTagsInput {
    splitChars: string[] | undefined;
    allowDuplicates: boolean | undefined;
    maxTags: number | undefined;
    value: string;
    currentTags: string[];
    isDuplicate?: (value: string, currentValues: string[]) => boolean;
    onDuplicate?: (value: string) => void;
}
export declare function getSplittedTags({ splitChars, allowDuplicates, maxTags, value, currentTags, isDuplicate, onDuplicate, }: GetSplittedTagsInput): string[];
export {};
