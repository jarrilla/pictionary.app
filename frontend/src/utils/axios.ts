import axios from 'axios'

// Get the API URL from environment variables - this is required in both dev and prod
const apiUrl = import.meta.env.VITE_API_URL

if (!apiUrl) {
  console.error('VITE_API_URL environment variable is not set')
}

// Create axios instance with API URL as baseURL
const instance = axios.create({
  baseURL: apiUrl,
  headers: {
    'Content-Type': 'application/json',
  },
})

export default instance 