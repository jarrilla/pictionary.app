<template>
  <div class="mobile-home">
    <div
      :class="{
        'search-container': true,
        'search-top': hasDefinition
      }"
      :style="{
        'height': hasDefinition ? '3.5rem' : '8rem'
      }"
    >
      <h1 v-if="!hasDefinition">The dictionary for visual learners!</h1>
      <div class="search-box">
        <input 
          type="text" 
          v-model="searchWord" 
          @keyup.enter="lookupWord"
          :placeholder="currentWord || 'Enter a word...'"
          :disabled="isLoading"
        />
        <button @click="lookupWord" :disabled="isLoading || !searchWord">
          {{ isLoading ? 'Loading...' : 'Search' }}
        </button>
      </div>
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <div v-if="hasDefinition" class="result-container">
      <div class="definition-section">

        <button
          @click="previousDefinition"
          :disabled="currentDefinitionIndex === 0"
          class="nav-button-previous"
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" class="bi bi-chevron-compact-left" viewBox="0 0 16 16">
            <path fill-rule="evenodd" d="M9.224 1.553a.5.5 0 0 1 .223.67L6.56 8l2.888 5.776a.5.5 0 1 1-.894.448l-3-6a.5.5 0 0 1 0-.448l3-6a.5.5 0 0 1 .67-.223"/>
          </svg>
        </button>

        <div class="current-definition" v-if="currentDefinitionData">
          <div class="definition-header">
            <div class="part-of-speech">{{ currentDefinitionData.partOfSpeech }}</div>
            <div class="definition-counter">
              {{ currentDefinitionIndex + 1 }} / {{ flatDefinitions.length }}
            </div>
          </div>
          <p class="definition-text">{{ currentDefinitionData.definition }}</p>
          <p v-if="currentDefinitionData.example" class="example">
            Example: "{{ currentDefinitionData.example }}"
          </p>
        </div>

        <button
          class="nav-button-next"
          @click="nextDefinition"
          :disabled="currentDefinitionIndex >= flatDefinitions.length - 1"
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="currentColor" class="bi bi-chevron-compact-right" viewBox="0 0 16 16">
            <path fill-rule="evenodd" d="M6.776 1.553a.5.5 0 0 1 .671.223l3 6a.5.5 0 0 1 0 .448l-3 6a.5.5 0 1 1-.894-.448L9.44 8 6.553 2.224a.5.5 0 0 1 .223-.671"/>
          </svg>
        </button>
      </div>

      <div class="image-section">
        <div v-if="isImageLoading" class="loading-spinner">
          <div class="spinner"></div>
          <p>Generating image...</p>
        </div>
        <div v-else-if="imageError" class="error-message">
          <p>{{ imageError }}</p>
        </div>
        <div v-else-if="imageUrl" class="image-display">
          <img :src="imageUrl" alt="Generated image for the word" />
          <button 
            class="regenerate-btn" 
            @click="regenerateImage"
            :disabled="isImageLoading"
            :title="'Click me if the image above wasn\'t helpful.'"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" class="bi bi-emoji-frown" viewBox="0 0 16 16">
              <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14m0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16"/>
              <path d="M4.285 12.433a.5.5 0 0 0 .683-.183A3.5 3.5 0 0 1 8 10.5c1.295 0 2.426.703 3.032 1.75a.5.5 0 0 0 .866-.5A4.5 4.5 0 0 0 8 9.5a4.5 4.5 0 0 0-3.898 2.25.5.5 0 0 0 .183.683M7 6.5C7 7.328 6.552 8 6 8s-1-.672-1-1.5S5.448 5 6 5s1 .672 1 1.5m4 0c0 .828-.448 1.5-1 1.5s-1-.672-1-1.5S9.448 5 10 5s1 .672 1 1.5"/>
            </svg>
            <span>Regenerate</span>
          </button>
        </div>
        <div v-else class="no-image">
          <p>No image available</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import axios from '@/utils/axios'

const searchWord = ref('')
const wordData = ref<any>(null)
const currentWord = ref('')
const isLoading = ref(false)
const error = ref('')
const imageUrl = ref('')
const isImageLoading = ref(false)
const imageError = ref('')
const currentDefinitionIndex = ref(0)
const currentImageController = ref<AbortController | null>(null)

const hasDefinition = computed(() => wordData.value !== null)

interface FlatDefinition {
  partOfSpeech: string;
  definition: string;
  example?: string;
}

const flatDefinitions = computed<FlatDefinition[]>(() => {
  if (!wordData.value) return []
  
  return wordData.value.meanings.flatMap((meaning: any) => 
    meaning.definitions.map((def: any) => ({
      partOfSpeech: meaning.partOfSpeech,
      definition: def.definition,
      example: def.example
    }))
  )
})

const currentDefinitionData = computed(() => {
  if (!flatDefinitions.value.length) return null
  return flatDefinitions.value[currentDefinitionIndex.value]
})

const nextDefinition = () => {
  if (currentDefinitionIndex.value < flatDefinitions.value.length - 1) {
    currentDefinitionIndex.value++
    generateImageForCurrentDefinition()
  }
}

const previousDefinition = () => {
  if (currentDefinitionIndex.value > 0) {
    currentDefinitionIndex.value--
    generateImageForCurrentDefinition()
  }
}

const lookupWord = async () => {
  if (!searchWord.value || isLoading.value) return

  isLoading.value = true
  error.value = ''
  wordData.value = null
  imageUrl.value = ''
  currentDefinitionIndex.value = 0
  
  try {
    const response = await axios.get(`https://api.dictionaryapi.dev/api/v2/entries/en/${searchWord.value.trim()}`)
    
    if (response.data && response.data.length > 0) {
      wordData.value = response.data[0]
      currentWord.value = searchWord.value.trim()
      generateImageForCurrentDefinition()
    }
  } catch (err: any) {
    if (err.response && err.response.status === 404) {
      error.value = 'Word not found. Please try another word.'
    } else {
      error.value = 'An error occurred while looking up the word. Please try again.'
    }
  } finally {
    isLoading.value = false
  }
}

const generateImageForCurrentDefinition = async () => {
  if (!currentDefinitionData.value) return
  
  const { partOfSpeech, definition } = currentDefinitionData.value
  generateImage(currentWord.value, partOfSpeech, definition)
}

const generateImage = async (word: string, partOfSpeech: string, definition: string) => {
  isImageLoading.value = true
  imageError.value = ''
  
  if (currentImageController.value) {
    currentImageController.value.abort()
  }

  currentImageController.value = new AbortController()
  
  try {
    const cacheResponse = await axios.get('/api/cache', {
      params: { word, partOfSpeech, definition }
    })
    
    if (cacheResponse.data && cacheResponse.data.imageData) {
      imageUrl.value = `data:image/png;base64,${cacheResponse.data.imageData}`
      isImageLoading.value = false
      return
    }

    await regenerateImage()
  } catch (err: any) {
    if (err.response && err.response.status !== 404) {
      console.error('Cache check failed:', err)
    }
    await regenerateImage()
  }
}

const regenerateImage = async () => {
  if (!currentDefinitionData.value) return

  imageUrl.value = ''
  isImageLoading.value = true
  imageError.value = ''

  try {
    const response = await axios.post('/api/generate-image', {
      word: currentWord.value,
      partOfSpeech: currentDefinitionData.value.partOfSpeech,
      definition: currentDefinitionData.value.definition
    }, {
      signal: currentImageController.value?.signal
    })
    
    if (!currentImageController.value?.signal.aborted) {
      if (response.data && response.data.imageData) {
        imageUrl.value = `data:image/png;base64,${response.data.imageData}`
      } else if (response.data && response.data.error) {
        imageError.value = response.data.error
      }
    }
  } catch (err: any) {
    if (!currentImageController.value?.signal.aborted) {
      if (err.response?.data?.error) {
        imageError.value = err.response.data.error
      } else {
        imageError.value = 'Failed to generate image. Please try again.'
      }
      console.error('Image generation failed:', err)
    }
  } finally {
    if (!currentImageController.value?.signal.aborted) {
      isImageLoading.value = false
    }
  }
}
</script>

<style scoped>
.mobile-home {
  display: flex;
  flex-direction: column;
  height: calc(100% + 2rem);
  box-sizing: border-box;
  margin: -1rem;
}

.search-container {
  width: 100%;
  transition: all 0.3s ease;
  box-sizing: border-box;
  display: flex;
  flex-flow: column nowrap;
  align-items: center;
  justify-content: center;
}

.search-container h1 {
  margin-bottom: 10px;
  color: #2c3e50;
  font-size: 1.5rem;
  text-align: center;
}

.search-box {
  display: flex;
  justify-content: center;
  width: 100%;
  padding: 0 1rem;
  box-sizing: border-box;
}

.search-box input {
  padding: 10px 15px;
  flex: 1;
  max-width: 300px;
  border: 1px solid #ddd;
  border-radius: 4px 0 0 4px;
  font-size: 16px;
}

.search-box button {
  padding: 10px 20px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 0 4px 4px 0;
  cursor: pointer;
  font-size: 16px;
  white-space: nowrap;
}

.search-box button:disabled {
  background-color: #cccccc;
}

.result-container {
  display: flex;
  flex-direction: column;
  height: calc(100% - 3.5rem);
  box-sizing: border-box;
}

.definition-section {
  background: white;
  padding: 0.25rem 1.25rem;
  position: relative;
  height: 40%;
  box-sizing: border-box;
  border-bottom: 1px solid #ddd;
  /* overflow-y: auto; */
}

.nav-button-previous,
.nav-button-next {
  position: absolute;
  top: 0;
  height: 100%;
  width: 1rem;
  z-index: 10;
  background-color: transparent;
  color: black;
  font-size: 1.2rem;
  border: none;
  user-select: none;
  outline: none;
  cursor: pointer;
}
.nav-button-previous:disabled,
.nav-button-next:disabled {
  opacity: 0;
  display: none;
}
.nav-button-previous {
  left: 0;
}
.nav-button-next {
  right: 0;
}

.nav-button-previous svg,
.nav-button-next svg {
  width: 100%;
  height: 100%;
  transform: scale(2);
}

.definition-counter {
  font-size: 0.9rem;
  color: #666;
}

.current-definition {
  font-size: 0.85rem;
  padding: 0.5rem;
  background: #f9f9f9;
}

.definition-header {
  display: flex;
  flex-flow: row nowrap;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  font-style: italic;
}

.part-of-speech {
  color: #42b983;
}

.definition-text {
  font-size: 0.85rem;
  /* line-height: 1.2; */
  margin-bottom: 0.5rem;
}

.example {
  color: #666;
  font-style: italic;
  font-size: 0.85rem;
}

.image-section {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 60%;
  box-sizing: border-box;
}

.image-display {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  box-sizing: border-box;
  height: 100%;
}

.image-display img {
  width: auto;
  height: calc(100% - 2rem);
  border-radius: 8px;
  box-sizing: border-box;
}

.regenerate-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  background: transparent;
  border: 2px solid #2c3e50;
  border-radius: 4px;
  color: #2c3e50;
  cursor: pointer;
  transition: all 0.2s;
  width: 8rem;
  justify-content: center;
  height: 1.5rem;
}

.regenerate-btn:hover {
  background: #2c3e50;
  color: white;
}

.regenerate-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-top-color: #42b983;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-message {
  color: #e74c3c;
  text-align: center;
  padding: 1rem;
  background-color: rgba(231, 76, 60, 0.1);
  border-radius: 4px;
  margin: 1rem;
}

.no-image {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  max-width: 300px;
  height: 200px;
  background: #f5f5f5;
  border-radius: 8px;
  color: #666;
}
</style> 