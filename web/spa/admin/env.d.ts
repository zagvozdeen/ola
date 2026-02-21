/// <reference types="vite/client" />
import 'telegram-web-app'

declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    const component: DefineComponent<object, object, unknown>
    export default component
}
