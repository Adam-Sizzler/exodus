export declare const CRUD_ACTIONS: {
    readonly CREATED: "CREATED";
    readonly UPDATED: "UPDATED";
    readonly DELETED: "DELETED";
};
export type TCrudActions = [keyof typeof CRUD_ACTIONS][number];
export type TCrudActionsKeys = (typeof CRUD_ACTIONS)[keyof typeof CRUD_ACTIONS];
export declare const CRUD_ACTIONS_VALUES: ("CREATED" | "UPDATED" | "DELETED")[];
//# sourceMappingURL=crud-actions.constant.d.ts.map