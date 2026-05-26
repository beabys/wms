/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** When true, the mock API layer intercepts fetch calls */
  readonly VITE_USE_MOCK?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
