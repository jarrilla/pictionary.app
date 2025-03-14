// This file is used to help TypeScript recognize Vue APIs
import { DefineComponent } from 'vue'

declare module 'vue' {
  export { ref, computed } from '@vue/runtime-core'
} 