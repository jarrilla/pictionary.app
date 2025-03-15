<template>
  <div class="home">
    <div
      :class="{
        'search-container': true,
        'search-top': hasDefinition
      }"
      :style="{
        'height': hasDefinition ? '5rem' : '10rem'
      }"
    >
      <h1 v-if="!hasDefinition">Type a word to look up its definition</h1>
      <div class="search-box">
        <input 
          type="text" 
          v-model="searchWord" 
          @keyup.enter="lookupWord"
          placeholder="Enter a word..."
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
      <div class="definition-container">
        <h2>{{ currentWord }}</h2>

        <div class="meanings-container">
          <div v-for="(meaning, index) in wordData.meanings" :key="index" class="meaning">
            <p class="part-of-speech">{{ meaning.partOfSpeech }}</p>
            <div
              v-for="(definition, defIndex) in meaning.definitions"
              :key="defIndex"
              class="definition"
              @click="selectDefinition(index, defIndex)"
              :class="{
                'selected': selectedDefinition.meaningIndex === index && selectedDefinition.definitionIndex === defIndex
              }"
            >
              <p>{{ definition.definition }}</p>
              <p v-if="definition.example" class="example">Example: "{{ definition.example }}"</p>
            </div>
          </div>
        </div>
        
      </div>
      <div class="image-container">
        <div v-if="isImageLoading" class="loading-spinner">
          <div class="spinner"></div>
          <p>Generating image...</p>
        </div>
        <div v-else-if="imageError" class="error-message">
          <p>{{ imageError }}</p>
        </div>
        <div v-else-if="imageUrl" class="image-display">
          <img :src="imageUrl" alt="Generated image for the word" />
        </div>
        <div v-else class="no-image">
          <p>No image available</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from '@vue/runtime-dom'
import axios from 'axios'

const searchWord = ref('')
const wordData = ref<any>(null)
const currentWord = ref('')
const isLoading = ref(false)
const error = ref('')
const imageUrl = ref('')
const isImageLoading = ref(false)
const imageError = ref('')
const selectedDefinition = ref({
  meaningIndex: 0,
  definitionIndex: 0
})

const hasDefinition = computed(() => wordData.value !== null)

const lookupWord = async () => {
  if (!searchWord.value || isLoading.value) return

  isLoading.value = true
  error.value = ''
  wordData.value = null
  imageUrl.value = ''
  
  try {
    // Fetch word definition from Dictionary API
    const response = await axios.get(`https://api.dictionaryapi.dev/api/v2/entries/en/${searchWord.value.trim()}`)
    
    if (response.data && response.data.length > 0) {
      wordData.value = response.data[0]
      currentWord.value = searchWord.value.trim()
      
      // Get the first definition to use in the image generation
      const firstMeaning = wordData.value.meanings[0]
      const partOfSpeech = firstMeaning.partOfSpeech
      const definition = firstMeaning.definitions[0].definition
      
      // Generate image using backend API
      generateImage(currentWord.value, partOfSpeech, definition)
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

const generateImage = async (word: string, partOfSpeech: string, definition: string) => {
  isImageLoading.value = true
  imageError.value = ''
  
  try {
    const response = await axios.post('/api/generate-image', {
      word,
      partOfSpeech,
      definition
    })
    
    if (response.data && response.data.imageUrl) {
      imageUrl.value = response.data.imageUrl
    } else if (response.data && response.data.error) {
      imageError.value = response.data.error
    }
  } catch (err: any) {
    if (err.response && err.response.data && err.response.data.error) {
      imageError.value = err.response.data.error
    } else {
      imageError.value = 'Failed to generate image. Please try again.'
    }
  } finally {
    isImageLoading.value = false
  }
}

const selectDefinition = (meaningIndex: number, definitionIndex: number) => {
  selectedDefinition.value.meaningIndex = meaningIndex
  selectedDefinition.value.definitionIndex = definitionIndex

  const partOfSpeech = wordData.value.meanings[meaningIndex].partOfSpeech
  const definition = wordData.value.meanings[meaningIndex].definitions[definitionIndex].definition

  generateImage(currentWord.value, partOfSpeech, definition)
}
</script>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100%;
  box-sizing: border-box;
}

.search-container {
  width: 100%;
  transition: all 0.5s ease;
  box-sizing: border-box;
  display: flex;
  flex-flow: column nowrap;
  align-items: center;
  justify-content: center;
  padding: 0.5rem;
}

.search-container h1 {
  margin-bottom: 20px;
  color: #2c3e50;
}

.search-box {
  display: flex;
  justify-content: center;
}

.search-box input {
  padding: 10px 15px;
  width: 300px;
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
  transition: background-color 0.3s;
}

.search-box button:hover {
  background-color: #3aa876;
}

.search-box button:disabled {
  background-color: #cccccc;
  cursor: not-allowed;
}

.error-message {
  color: #e74c3c;
  margin: 20px 0;
  text-align: center;
  padding: 10px;
  border-radius: 4px;
  background-color: rgba(231, 76, 60, 0.1);
  border: 1px solid rgba(231, 76, 60, 0.3);
}

.image-container .error-message {
  width: 300px;
  height: auto;
  min-height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.result-container {
  display: flex;
  width: 100%;
  border-top: 1px solid #ddd;
  height: calc(100% - 5rem);
  box-sizing: border-box;
}

.definition-container {
  flex: 1;
  padding-right: 0.5rem;
  height: 100%;
  box-sizing: border-box;
}

.definition-container h2 {
  font-size: 28px;
  color: #2c3e50;
  margin-bottom: 0.5rem;
  text-transform: capitalize;
  height: 3rem;
}

.meanings-container {
  height: calc(100% - 3rem - 0.5rem);
  overflow-y: auto;
  scrollbar-width: thin;
  box-sizing: border-box;
  display: flex;
  flex-flow: column nowrap;
}

.meaning {
  margin-bottom: 20px;
}

.part-of-speech {
  font-style: italic;
  color: #7f8c8d;
  margin-bottom: 10px;
}

.definition {
  margin-bottom: 15px;
  padding-left: 15px;
  border-left: 3px solid #42b983;
  cursor: pointer;
  transition: background-color 0.25s;
}

.definition:hover,
.definition.selected {
  background-color: #42b98335;
}

.example {
  color: #7f8c8d;
  font-style: italic;
  margin-top: 5px;
}

.image-container {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  width: 100%;
}

.image-container img {
  max-width: 100%;
  border-radius: 8px;
}

.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 300px;
  height: 300px;
  background-color: #f5f5f5;
  border-radius: 8px;
  color: #7f8c8d;
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
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

.image-display {
  width: 300px;
  height: 300px;
  border-radius: 8px;
  overflow: hidden;
  margin-top: 20px;
}

.no-image {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 300px;
  height: 300px;
  background-color: #f5f5f5;
  border-radius: 8px;
  color: #7f8c8d;
}

@media (max-width: 768px) {
  .result-container {
    flex-direction: column;
  }
  
  .definition-container {
    padding-right: 0;
    margin-bottom: 20px;
  }
  
  .search-box input {
    width: 200px;
  }
}
</style> 