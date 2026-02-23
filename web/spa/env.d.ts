/// <reference types="vite/client" />
import 'telegram-web-app'

interface ImportMetaEnv {
    readonly VITE_API_URL: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}

declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    const component: DefineComponent<object, object, unknown>
    export default component
}
