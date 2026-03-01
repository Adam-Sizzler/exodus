const DEFAULT_TIMEOUT_MS = 4500;
const VALIDATE_METHODS = ['validate', 'validateConfig', 'check', 'checkConfig', 'lint'];
const FORMAT_METHODS = ['format', 'formatConfig', 'prettify', 'normalize'];

const MODULE_DEFS = {
  xray: {
    runtimeUrl: import.meta.env.VITE_GO_WASM_RUNTIME_URL || '/wasm/wasm_exec.js',
    wasmUrl: import.meta.env.VITE_XRAY_WASM_URL || '/wasm/xray.wasm',
    apiGlobal: import.meta.env.VITE_XRAY_WASM_GLOBAL || 'xrayWasm',
  },
  mihomo: {
    runtimeUrl: import.meta.env.VITE_GO_WASM_RUNTIME_URL || '/wasm/wasm_exec.js',
    wasmUrl: import.meta.env.VITE_MIHOMO_WASM_URL || '/wasm/mihomo.wasm',
    apiGlobal: import.meta.env.VITE_MIHOMO_WASM_GLOBAL || 'mihomoWasm',
  },
};

const runtimePromises = new Map();
const modulePromises = new Map();

const wait = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

const ensureScript = (url) => {
  if (typeof globalThis.Go === 'function') {
    return Promise.resolve();
  }

  const existing = runtimePromises.get(url);
  if (existing) {
    return existing;
  }

  const promise = new Promise((resolve, reject) => {
    const script = document.createElement('script');
    script.src = url;
    script.async = true;
    script.onload = () => {
      if (typeof globalThis.Go !== 'function') {
        reject(new Error(`WASM runtime loaded but Go constructor is missing: ${url}`));
        return;
      }
      resolve();
    };
    script.onerror = () => reject(new Error(`Failed to load WASM runtime: ${url}`));
    document.head.appendChild(script);
  });

  runtimePromises.set(url, promise);
  return promise;
};

const instantiateGoModule = async (wasmUrl, importObject) => {
  if ('instantiateStreaming' in WebAssembly) {
    try {
      const streamingResult = await WebAssembly.instantiateStreaming(fetch(wasmUrl), importObject);
      return streamingResult.instance || streamingResult;
    } catch (err) {
      // Fallback for servers with wrong WASM MIME type.
      if (!String(err?.message || '').toLowerCase().includes('mime')) {
        throw err;
      }
    }
  }

  const response = await fetch(wasmUrl);
  if (!response.ok) {
    throw new Error(`Failed to fetch WASM module: ${wasmUrl} (${response.status})`);
  }
  const bytes = await response.arrayBuffer();
  const moduleResult = await WebAssembly.instantiate(bytes, importObject);
  return moduleResult.instance || moduleResult;
};

const resolveApiFromGlobal = (engine, explicitGlobalName) => {
  const names = [
    explicitGlobalName,
    `${engine}Wasm`,
    `${engine}API`,
    `${engine}Api`,
    engine,
    `${engine}Validator`,
    `${engine}Bridge`,
  ].filter(Boolean);

  for (const name of names) {
    const value = globalThis[name];
    if (value && (typeof value === 'object' || typeof value === 'function')) {
      return { api: value, name };
    }
  }

  const directValidateName = `${engine}Validate`;
  const directFormatName = `${engine}Format`;
  const validateFn = globalThis[directValidateName];
  const formatFn = globalThis[directFormatName];

  if (typeof validateFn === 'function') {
    return {
      api: {
        validate: validateFn.bind(globalThis),
        format: typeof formatFn === 'function' ? formatFn.bind(globalThis) : undefined,
      },
      name: `${directValidateName}${typeof formatFn === 'function' ? `/${directFormatName}` : ''}`,
    };
  }

  return null;
};

const waitForApi = async (engine, globalName, timeoutMs) => {
  const started = Date.now();
  while (Date.now() - started < timeoutMs) {
    const discovered = resolveApiFromGlobal(engine, globalName);
    if (discovered) {
      return discovered;
    }
    await wait(25);
  }
  throw new Error(
    `WASM API was not initialized in time for "${engine}". Expected global like "${globalName}" or "${engine}Wasm".`
  );
};

const normalizeValidationResult = (result) => {
  if (result && typeof result === 'object' && typeof result.ok === 'boolean') {
    return {
      ok: result.ok,
      message: typeof result.message === 'string' ? result.message : (result.ok ? 'Config is valid.' : 'Config is invalid.'),
      formatted: typeof result.formatted === 'string' ? result.formatted : undefined,
    };
  }

  if (typeof result === 'boolean') {
    return {
      ok: result,
      message: result ? 'Config is valid.' : 'Config is invalid.',
    };
  }

  if (typeof result === 'string') {
    return {
      ok: true,
      message: 'Config is valid.',
      formatted: result,
    };
  }

  return {
    ok: true,
    message: 'Config is valid.',
  };
};

const resolveCallable = (api, methodNames) => {
  for (const methodName of methodNames) {
    if (api && typeof api === 'object' && typeof api[methodName] === 'function') {
      return api[methodName].bind(api);
    }
  }
  return null;
};

const loadModule = async (engine, timeoutMs = DEFAULT_TIMEOUT_MS) => {
  const def = MODULE_DEFS[engine];
  if (!def) {
    throw new Error(`Unknown WASM engine: ${engine}`);
  }

  const cached = modulePromises.get(engine);
  if (cached) {
    return cached;
  }

  const promise = (async () => {
    await ensureScript(def.runtimeUrl);

    const go = new globalThis.Go();
    const instance = await instantiateGoModule(def.wasmUrl, go.importObject);

    // Go WASM programs are usually long-running and expose API via global object.
    go.run(instance).catch((err) => {
      console.error(`[wasm:${engine}] runtime exited with error`, err);
    });

    const discovered = await waitForApi(engine, def.apiGlobal, timeoutMs);
    const api = discovered.api;
    const validate = resolveCallable(api, VALIDATE_METHODS);
    const format = resolveCallable(api, FORMAT_METHODS);

    if (!validate) {
      throw new Error(
        `WASM API "${discovered.name}" does not expose any validate method (${VALIDATE_METHODS.join(', ')})`
      );
    }

    return { validate, format };
  })();

  modulePromises.set(engine, promise);
  return promise;
};

export const validateWithWasm = async (engine, input, timeoutMs = DEFAULT_TIMEOUT_MS) => {
  const moduleApi = await loadModule(engine, timeoutMs);
  const result = await moduleApi.validate(input);
  return normalizeValidationResult(result);
};

export const formatWithWasm = async (engine, input, timeoutMs = DEFAULT_TIMEOUT_MS) => {
  const moduleApi = await loadModule(engine, timeoutMs);
  if (!moduleApi.format) {
    throw new Error(`WASM engine "${engine}" does not expose format(text)`);
  }
  const result = await moduleApi.format(input);
  if (typeof result === 'string') {
    return result;
  }
  if (result && typeof result === 'object' && typeof result.formatted === 'string') {
    return result.formatted;
  }
  throw new Error(`WASM format() returned unsupported result for engine "${engine}"`);
};
