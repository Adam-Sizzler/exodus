import { QueryFunction, MutateFunction } from '@tanstack/query-core';
import { UseQueryOptions, QueryFunction as QueryFunction$1 } from '@tanstack/react-query';

/**
 * @internal
 */
declare type InternalKey = `_${string}`;
/**
 * @internal
 */
declare type ExtractInternalKeys<T extends Record<string, unknown>> = Extract<keyof T, InternalKey>;

declare type AnyMutableOrReadonlyArray = any[] | readonly any[];
declare type Tuple = [ValidValue | undefined, ...Array<ValidValue | undefined>];
declare type KeyTuple = Tuple | Readonly<Tuple>;
declare type ValidValue = string | number | boolean | object | bigint;
declare type Prettify<T> = {
    [K in keyof T]: T[K];
} & {};
declare type DefinitionKey<Key extends AnyMutableOrReadonlyArray> = {
    _def: readonly [...Key];
};

declare type ValidFactoryKey = 'queryKey' | 'queryFn' | 'contextQueries';
declare type StrictOptions<T> = T extends any[] ? T : keyof T extends ValidFactoryKey ? T : {
    [K in keyof T]: K extends ValidFactoryKey ? T[K] : never;
};
declare type NullableQueryKeyRecord = Record<'queryKey', KeyTuple | null>;
declare type QueryKeyRecord = Record<'queryKey', KeyTuple>;
declare type KeySchemaWithContextualQueries = NullableQueryKeyRecord & {
    contextQueries: QueryFactorySchema;
};
declare type $QueryFactorySchema = NullableQueryKeyRecord & {
    queryFn: QueryFunction;
};
declare type QueryFactoryWithContextualQueriesSchema = NullableQueryKeyRecord & {
    queryFn: QueryFunction;
    contextQueries: QueryFactorySchema;
};
declare type DynamicKeySchemaWithContextualQueries = QueryKeyRecord & {
    contextQueries: QueryFactorySchema;
};
declare type DynamicQueryFactorySchema = QueryKeyRecord & {
    queryFn: QueryFunction;
};
declare type DynamicQueryFactoryWithContextualQueriesSchema = QueryKeyRecord & {
    queryFn: QueryFunction;
    contextQueries: QueryFactorySchema;
};
declare type FactoryProperty = null | KeyTuple | NullableQueryKeyRecord | KeySchemaWithContextualQueries | $QueryFactorySchema | QueryFactoryWithContextualQueriesSchema;
declare type DynamicKey = (...args: any[]) => DynamicQueryFactoryWithContextualQueriesSchema | DynamicQueryFactorySchema | DynamicKeySchemaWithContextualQueries | QueryKeyRecord | KeyTuple;
declare type QueryFactorySchema = Record<string, FactoryProperty | DynamicKey>;
declare type InvalidSchema$1<Schema extends QueryFactorySchema> = Omit<Schema, InternalKey>;
declare type ValidateFactory$1<Schema extends QueryFactorySchema> = ExtractInternalKeys<Schema> extends never ? {
    [P in keyof Schema]: Schema[P] extends (...args: infer Args) => infer R ? (...args: Args) => StrictOptions<R> : Schema[P];
} : InvalidSchema$1<Schema>;
declare type ExtractNullableKey$1<Key extends KeyTuple | null | undefined> = Key extends [...infer Value] | readonly [...infer Value] ? Value : Key extends null | undefined | unknown ? null : never;
declare type ComposeQueryKey$1<BaseKey extends AnyMutableOrReadonlyArray, Key> = Key extends KeyTuple ? readonly [...BaseKey, ...Key] : readonly [...BaseKey];
declare type QueryOptionsStruct<Keys extends AnyMutableOrReadonlyArray, Fetcher extends QueryFunction, FetcherResult extends ReturnType<Fetcher> = ReturnType<Fetcher>> = {
    queryKey: readonly [...Keys];
    queryFn: QueryFunction<Awaited<FetcherResult>, readonly [...Keys]>;
};
declare type FactoryWithContextualQueriesOutput<BaseKey extends AnyMutableOrReadonlyArray, Schema extends KeySchemaWithContextualQueries | DynamicKeySchemaWithContextualQueries, SchemaQueryKey extends Schema['queryKey'] = Schema['queryKey'], ContextQueries extends Schema['contextQueries'] = Schema['contextQueries'], ComposedKey extends AnyMutableOrReadonlyArray = ComposeQueryKey$1<BaseKey, ExtractNullableKey$1<SchemaQueryKey>>> = SchemaQueryKey extends null ? Omit<QueryOptionsStruct<ComposedKey, QueryFunction>, 'queryFn'> & {
    _ctx: {
        [P in keyof ContextQueries]: ContextQueries[P] extends DynamicKey ? DynamicFactoryOutput<[...ComposedKey, P], ContextQueries[P]> : ContextQueries[P] extends FactoryProperty ? StaticFactoryOutput<[...ComposedKey, P], ContextQueries[P]> : never;
    };
} : Omit<QueryOptionsStruct<ComposedKey, QueryFunction>, 'queryFn'> & DefinitionKey<BaseKey> & {
    _ctx: {
        [P in keyof ContextQueries]: ContextQueries[P] extends DynamicKey ? DynamicFactoryOutput<[...ComposedKey, P], ContextQueries[P]> : ContextQueries[P] extends FactoryProperty ? StaticFactoryOutput<[...ComposedKey, P], ContextQueries[P]> : never;
    };
};
declare type FactoryQueryKeyRecordOutput<BaseKey extends AnyMutableOrReadonlyArray, Schema extends NullableQueryKeyRecord | QueryKeyRecord, SchemaQueryKey extends Schema['queryKey'] = Schema['queryKey'], ComposedKey extends AnyMutableOrReadonlyArray = ComposeQueryKey$1<BaseKey, ExtractNullableKey$1<SchemaQueryKey>>> = SchemaQueryKey extends null ? Omit<QueryOptionsStruct<BaseKey, QueryFunction>, 'queryFn'> : Omit<QueryOptionsStruct<ComposedKey, QueryFunction>, 'queryFn'> & DefinitionKey<BaseKey>;
declare type FactoryQueryOptionsOutput<BaseKey extends AnyMutableOrReadonlyArray, Schema extends $QueryFactorySchema | DynamicQueryFactorySchema, SchemaQueryKey extends Schema['queryKey'] = Schema['queryKey'], QueryFn extends Schema['queryFn'] = Schema['queryFn'], ComposedKey extends AnyMutableOrReadonlyArray = ComposeQueryKey$1<BaseKey, ExtractNullableKey$1<SchemaQueryKey>>> = SchemaQueryKey extends null ? QueryOptionsStruct<BaseKey, QueryFn> : QueryOptionsStruct<ComposedKey, QueryFn> & DefinitionKey<BaseKey>;
declare type FactoryQueryOptionsWithContextualQueriesOutput<BaseKey extends AnyMutableOrReadonlyArray, Schema extends QueryFactoryWithContextualQueriesSchema | DynamicQueryFactoryWithContextualQueriesSchema, SchemaQueryKey extends Schema['queryKey'] = Schema['queryKey'], QueryFn extends Schema['queryFn'] = Schema['queryFn'], ContextQueries extends Schema['contextQueries'] = Schema['contextQueries'], Key extends AnyMutableOrReadonlyArray = ComposeQueryKey$1<BaseKey, ExtractNullableKey$1<SchemaQueryKey>>> = SchemaQueryKey extends null ? QueryOptionsStruct<Key, QueryFn> & {
    _ctx: {
        [P in keyof ContextQueries]: ContextQueries[P] extends DynamicKey ? DynamicFactoryOutput<[...Key, P], ContextQueries[P]> : ContextQueries[P] extends FactoryProperty ? StaticFactoryOutput<[...Key, P], ContextQueries[P]> : never;
    };
} : DefinitionKey<BaseKey> & QueryOptionsStruct<Key, QueryFn> & {
    _ctx: {
        [P in keyof ContextQueries]: ContextQueries[P] extends DynamicKey ? DynamicFactoryOutput<[...Key, P], ContextQueries[P]> : ContextQueries[P] extends FactoryProperty ? StaticFactoryOutput<[...Key, P], ContextQueries[P]> : never;
    };
};
declare type DynamicFactoryOutput<Keys extends AnyMutableOrReadonlyArray, Generator extends DynamicKey, Output extends ReturnType<Generator> = ReturnType<Generator>> = ((...args: Parameters<Generator>) => Output extends [...infer TupleResult] | readonly [...infer TupleResult] ? Omit<QueryOptionsStruct<[...Keys, ...TupleResult], QueryFunction>, 'queryFn'> : Output extends DynamicQueryFactoryWithContextualQueriesSchema ? Omit<FactoryQueryOptionsWithContextualQueriesOutput<Keys, Output>, '_def'> : Output extends DynamicQueryFactorySchema ? Omit<FactoryQueryOptionsOutput<Keys, Output>, '_def'> : Output extends DynamicKeySchemaWithContextualQueries ? Omit<FactoryWithContextualQueriesOutput<Keys, Output>, '_def'> : Output extends QueryKeyRecord ? Omit<FactoryQueryKeyRecordOutput<Keys, Output>, '_def'> : never) & DefinitionKey<Keys>;
declare type AnyQueryFactoryOutputCallback = DynamicFactoryOutput<[string, ...any[]], DynamicKey>;
declare type StaticFactoryOutput<Keys extends AnyMutableOrReadonlyArray, Property extends FactoryProperty> = Property extends null ? Omit<QueryOptionsStruct<Keys, QueryFunction>, 'queryFn'> : Property extends [...infer Result] | readonly [...infer Result] ? DefinitionKey<Keys> & Omit<QueryOptionsStruct<[...Keys, ...Result], QueryFunction>, 'queryFn'> : Property extends QueryFactoryWithContextualQueriesSchema ? FactoryQueryOptionsWithContextualQueriesOutput<Keys, Property> : Property extends $QueryFactorySchema ? FactoryQueryOptionsOutput<Keys, Property> : Property extends KeySchemaWithContextualQueries ? FactoryWithContextualQueriesOutput<Keys, Property> : Property extends NullableQueryKeyRecord ? FactoryQueryKeyRecordOutput<Keys, Property> : never;
declare type FactoryOutput<Key extends string, Schema extends QueryFactorySchema> = DefinitionKey<[Key]> & {
    [P in keyof Schema]: Schema[P] extends DynamicKey ? DynamicFactoryOutput<[Key, P], Schema[P]> : Schema[P] extends FactoryProperty ? StaticFactoryOutput<[Key, P], Schema[P]> : never;
};
declare type QueryKeyFactoryResult<Key extends string, Schema extends QueryFactorySchema> = FactoryOutput<Key, Schema>;
declare type AnyQueryKeyFactoryResult = DefinitionKey<[string]> | QueryKeyFactoryResult<string, any>;

declare type QueryKeyStoreSchema = Record<string, null | QueryFactorySchema>;
declare type QueryKeyStore<StoreSchema extends QueryKeyStoreSchema> = {
    [P in keyof StoreSchema & string]: StoreSchema[P] extends QueryFactorySchema ? QueryKeyFactoryResult<P, StoreSchema[P]> : DefinitionKey<[P]>;
};
declare function createQueryKeyStore<StoreSchema extends QueryKeyStoreSchema>(schema: StoreSchema): QueryKeyStore<StoreSchema>;

declare type NullableMutationKeyRecord = Record<'mutationKey', KeyTuple | null>;
declare type MutationKeyRecord = Record<'mutationKey', KeyTuple>;
declare type MutationKeySchemaWithContextualMutations = NullableMutationKeyRecord & {
    contextMutations: MutationFactorySchema;
};
declare type $MutationFactorySchema = NullableMutationKeyRecord & {
    mutationFn: MutateFunction;
};
declare type MutationFactoryWithContextualMutationsSchema = NullableMutationKeyRecord & {
    mutationFn: MutateFunction;
    contextMutations: MutationFactorySchema;
};
declare type DynamicMutationKeySchemaWithContextualMutations = MutationKeyRecord & {
    contextMutations: MutationFactorySchema;
};
declare type DynamicMutationFactorySchema = MutationKeyRecord & {
    mutationFn: MutateFunction;
};
declare type DynamicMutationFactoryWithContextualMutationsSchema = MutationKeyRecord & {
    mutationFn: MutateFunction;
    contextMutations: MutationFactorySchema;
};
declare type MutationFactoryProperty = null | KeyTuple | NullableMutationKeyRecord | MutationKeySchemaWithContextualMutations | $MutationFactorySchema | MutationFactoryWithContextualMutationsSchema;
declare type MutationDynamicKey = (...args: any[]) => DynamicMutationFactoryWithContextualMutationsSchema | DynamicMutationFactorySchema | DynamicMutationKeySchemaWithContextualMutations | MutationKeyRecord | KeyTuple;
declare type MutationFactorySchema = Record<string, MutationFactoryProperty | MutationDynamicKey>;
declare type InvalidSchema<Schema extends MutationFactorySchema> = Omit<Schema, InternalKey>;
declare type ValidateFactory<Schema extends MutationFactorySchema> = Schema extends {
    [P in ExtractInternalKeys<Schema>]: Schema[P];
} ? InvalidSchema<Schema> : Schema;
declare type ExtractNullableKey<Key extends KeyTuple | null | undefined> = Key extends [...infer Value] | readonly [...infer Value] ? Value : Key extends null | undefined | unknown ? null : never;
declare type ComposeQueryKey<BaseKey extends AnyMutableOrReadonlyArray, Key> = Key extends KeyTuple ? readonly [...BaseKey, ...Key] : readonly [...BaseKey];
declare type MutationOptionsStruct<Keys extends AnyMutableOrReadonlyArray, Fetcher extends MutateFunction, FetcherResult extends ReturnType<Fetcher> = ReturnType<Fetcher>, FetcherVariables extends Parameters<Fetcher>[0] = Parameters<Fetcher>[0]> = {
    mutationKey: readonly [...Keys];
    mutationFn: MutateFunction<Awaited<FetcherResult>, unknown, FetcherVariables>;
};
declare type MutationFactoryWithContextualQueriesOutput<BaseKey extends AnyMutableOrReadonlyArray, Schema extends MutationKeySchemaWithContextualMutations | DynamicMutationKeySchemaWithContextualMutations, SchemaMutationKey extends Schema['mutationKey'] = Schema['mutationKey'], ContextMutations extends Schema['contextMutations'] = Schema['contextMutations'], ComposedKey extends AnyMutableOrReadonlyArray = ComposeQueryKey<BaseKey, ExtractNullableKey<SchemaMutationKey>>> = SchemaMutationKey extends null ? Omit<MutationOptionsStruct<ComposedKey, MutateFunction>, 'mutationFn'> & {
    _ctx: {
        [P in keyof ContextMutations]: ContextMutations[P] extends MutationDynamicKey ? DynamicMutationFactoryOutput<[...ComposedKey, P], ContextMutations[P]> : ContextMutations[P] extends MutationFactoryProperty ? StaticMutationFactoryOutput<[...ComposedKey, P], ContextMutations[P]> : never;
    };
} : Omit<MutationOptionsStruct<ComposedKey, MutateFunction>, 'mutationFn'> & DefinitionKey<BaseKey> & {
    _ctx: {
        [P in keyof ContextMutations]: ContextMutations[P] extends MutationDynamicKey ? DynamicMutationFactoryOutput<[...ComposedKey, P], ContextMutations[P]> : ContextMutations[P] extends MutationFactoryProperty ? StaticMutationFactoryOutput<[...ComposedKey, P], ContextMutations[P]> : never;
    };
};
declare type FactoryMutationKeyRecordOutput<BaseKey extends AnyMutableOrReadonlyArray, Schema extends NullableMutationKeyRecord | MutationKeyRecord, SchemaMutationKey extends Schema['mutationKey'] = Schema['mutationKey'], ComposedKey extends AnyMutableOrReadonlyArray = ComposeQueryKey<BaseKey, ExtractNullableKey<SchemaMutationKey>>> = SchemaMutationKey extends null ? Omit<MutationOptionsStruct<BaseKey, MutateFunction>, 'mutationFn'> : Omit<MutationOptionsStruct<ComposedKey, MutateFunction>, 'mutationFn'> & DefinitionKey<BaseKey>;
declare type FactoryMutationOptionsOutput<BaseKey extends AnyMutableOrReadonlyArray, Schema extends $MutationFactorySchema | DynamicMutationFactorySchema, SchemaQueryKey extends Schema['mutationKey'] = Schema['mutationKey'], MutationFn extends Schema['mutationFn'] = Schema['mutationFn'], ComposedKey extends AnyMutableOrReadonlyArray = ComposeQueryKey<BaseKey, ExtractNullableKey<SchemaQueryKey>>> = SchemaQueryKey extends null ? MutationOptionsStruct<BaseKey, MutationFn> : MutationOptionsStruct<ComposedKey, MutationFn> & DefinitionKey<BaseKey>;
declare type FactoryMutationOptionsWithContextualQueriesOutput<BaseKey extends AnyMutableOrReadonlyArray, Schema extends MutationFactoryWithContextualMutationsSchema | DynamicMutationFactoryWithContextualMutationsSchema, SchemaQueryKey extends Schema['mutationKey'] = Schema['mutationKey'], MutationFn extends Schema['mutationFn'] = Schema['mutationFn'], ContextMutations extends Schema['contextMutations'] = Schema['contextMutations'], Key extends AnyMutableOrReadonlyArray = ComposeQueryKey<BaseKey, ExtractNullableKey<SchemaQueryKey>>> = SchemaQueryKey extends null ? MutationOptionsStruct<Key, MutationFn> & {
    _ctx: {
        [P in keyof ContextMutations]: ContextMutations[P] extends MutationDynamicKey ? DynamicMutationFactoryOutput<[...Key, P], ContextMutations[P]> : ContextMutations[P] extends MutationFactoryProperty ? StaticMutationFactoryOutput<[...Key, P], ContextMutations[P]> : never;
    };
} : DefinitionKey<BaseKey> & MutationOptionsStruct<Key, MutationFn> & {
    _ctx: {
        [P in keyof ContextMutations]: ContextMutations[P] extends MutationDynamicKey ? DynamicMutationFactoryOutput<[...Key, P], ContextMutations[P]> : ContextMutations[P] extends MutationFactoryProperty ? StaticMutationFactoryOutput<[...Key, P], ContextMutations[P]> : never;
    };
};
declare type DynamicMutationFactoryOutput<Keys extends AnyMutableOrReadonlyArray, Generator extends MutationDynamicKey, Output extends ReturnType<Generator> = ReturnType<Generator>> = ((...args: Parameters<Generator>) => Output extends [...infer TupleResult] | readonly [...infer TupleResult] ? Omit<MutationOptionsStruct<[...Keys, ...TupleResult], MutateFunction>, 'mutationFn'> : Output extends DynamicMutationFactoryWithContextualMutationsSchema ? Omit<FactoryMutationOptionsWithContextualQueriesOutput<Keys, Output>, '_def'> : Output extends DynamicMutationFactorySchema ? Omit<FactoryMutationOptionsOutput<Keys, Output>, '_def'> : Output extends DynamicMutationKeySchemaWithContextualMutations ? Omit<MutationFactoryWithContextualQueriesOutput<Keys, Output>, '_def'> : Output extends MutationKeyRecord ? Omit<FactoryMutationKeyRecordOutput<Keys, Output>, '_def'> : never) & DefinitionKey<Keys>;
declare type AnyMutationFactoryOutputCallback = DynamicMutationFactoryOutput<[string, ...any[]], MutationDynamicKey>;
declare type StaticMutationFactoryOutput<Keys extends AnyMutableOrReadonlyArray, Property extends MutationFactoryProperty> = Property extends null ? Omit<MutationOptionsStruct<Keys, MutateFunction>, 'mutationFn'> : Property extends [...infer Result] | readonly [...infer Result] ? DefinitionKey<Keys> & Omit<MutationOptionsStruct<[...Keys, ...Result], MutateFunction>, 'mutationFn'> : Property extends MutationFactoryWithContextualMutationsSchema ? FactoryMutationOptionsWithContextualQueriesOutput<Keys, Property> : Property extends $MutationFactorySchema ? FactoryMutationOptionsOutput<Keys, Property> : Property extends MutationKeySchemaWithContextualMutations ? MutationFactoryWithContextualQueriesOutput<Keys, Property> : Property extends NullableMutationKeyRecord ? FactoryMutationKeyRecordOutput<Keys, Property> : never;
declare type MutationFactoryOutput<Key extends string, Schema extends MutationFactorySchema> = DefinitionKey<[Key]> & {
    [P in keyof Schema]: Schema[P] extends MutationDynamicKey ? DynamicMutationFactoryOutput<[Key, P], Schema[P]> : Schema[P] extends MutationFactoryProperty ? StaticMutationFactoryOutput<[Key, P], Schema[P]> : never;
};
declare type MutationKeyFactoryResult<Key extends string, Schema extends MutationFactorySchema> = MutationFactoryOutput<Key, Schema>;
declare type AnyMutationKeyFactoryResult = DefinitionKey<[string]> | MutationKeyFactoryResult<string, any>;

/**
 * @deprecated the type inference for this function is broken and will be fixed in the next patch version
 * or possibly removed and implemented differently in a major version
 */
declare function createMutationKeys<Key extends string>(mutationDef: Key): DefinitionKey<[Key]>;
/**
 * @deprecated the type inference for this function is broken and will be fixed in the next patch version
 * or possibly removed and implemented differently in a major version
 */
declare function createMutationKeys<Key extends string, Schema extends MutationFactorySchema>(mutationDef: Key, schema: ValidateFactory<Schema>): MutationKeyFactoryResult<Key, Schema>;

declare function createQueryKeys<Key extends string>(queryDef: Key): DefinitionKey<[Key]>;
declare function createQueryKeys<Key extends string, Schema extends QueryFactorySchema>(queryDef: Key, schema: ValidateFactory$1<Schema>): QueryKeyFactoryResult<Key, Schema>;

declare type StoreFromMergedQueryKeys<QueryOrMutationKeyFactoryResults extends Array<AnyQueryKeyFactoryResult | AnyMutationKeyFactoryResult>> = QueryOrMutationKeyFactoryResults extends ([
    infer First extends AnyQueryKeyFactoryResult | AnyMutationKeyFactoryResult,
    ...infer Rest extends AnyQueryKeyFactoryResult[] | AnyMutationKeyFactoryResult[]
]) ? {
    [P in First['_def'][0]]: First;
} & StoreFromMergedQueryKeys<Rest> : {};
declare function mergeQueryKeys<QueryKeyFactoryResults extends Array<AnyQueryKeyFactoryResult | AnyMutationKeyFactoryResult>>(...schemas: QueryKeyFactoryResults): Prettify<StoreFromMergedQueryKeys<QueryKeyFactoryResults>>;

declare type MergeInsertions<T> = T extends object ? {
    [K in keyof T]: MergeInsertions<T[K]>;
} : T;
declare type inferRecordMutationKeys<Target extends object> = {
    [P in Exclude<keyof Target, 'mutationFn'>]: Target[P] extends AnyMutableOrReadonlyArray ? Target[P] : Target[P] extends object ? {
        [K in keyof Target[P]]: inferSchemaProperty<Target[P][K]>;
    } : never;
};
declare type inferRecordQueryKeys<Target extends object> = {
    [P in Exclude<keyof Target, 'queryFn'>]: Target[P] extends AnyMutableOrReadonlyArray ? Target[P] : Target[P] extends object ? {
        [K in keyof Target[P]]: inferSchemaProperty<Target[P][K]>;
    } : never;
};
declare type inferSchemaProperty<Value> = Value extends AnyMutableOrReadonlyArray ? Value : Value extends StaticFactoryOutput<any[], any> ? inferRecordQueryKeys<Value> : Value extends StaticMutationFactoryOutput<any[], any> ? inferRecordMutationKeys<Value> : Value extends AnyQueryFactoryOutputCallback ? Record<'_def', Value['_def']> & inferRecordQueryKeys<ReturnType<Value>> : Value extends AnyMutationFactoryOutputCallback ? Record<'_def', Value['_def']> & inferRecordMutationKeys<ReturnType<Value>> : never;
declare type inferQueryKeys<Schema extends AnyQueryKeyFactoryResult | AnyMutationKeyFactoryResult> = {
    [P in keyof Schema]: MergeInsertions<inferSchemaProperty<Schema[P]>>;
};
declare type inferQueryKeyStore<Store extends QueryKeyStore<any>> = {
    [P in keyof Store]: inferQueryKeys<Store[P]>;
};
declare type LooseQueryOptionsStruct = {
    queryKey: AnyMutableOrReadonlyArray;
    queryFn: QueryFunction$1<any, any>;
};
declare type LooseQueryOptionsStructGenerator = (...args: any[]) => LooseQueryOptionsStruct;
declare type TypedUseQueryOptions<Options extends LooseQueryOptionsStruct | LooseQueryOptionsStructGenerator, Data = Options extends LooseQueryOptionsStructGenerator ? Awaited<ReturnType<ReturnType<Options>['queryFn']>> : Options extends LooseQueryOptionsStruct ? Awaited<ReturnType<Options['queryFn']>> : never> = Options extends LooseQueryOptionsStructGenerator ? UseQueryOptions<Awaited<ReturnType<ReturnType<Options>['queryFn']>>, unknown, Data, ReturnType<Options>['queryKey']> : Options extends LooseQueryOptionsStruct ? UseQueryOptions<Awaited<ReturnType<Options['queryFn']>>, unknown, Data, Options['queryKey']> : never;

export { type TypedUseQueryOptions, createMutationKeys, createQueryKeyStore, createQueryKeys, type inferQueryKeyStore, type inferQueryKeys, mergeQueryKeys };
