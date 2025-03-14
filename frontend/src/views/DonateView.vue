<template>
  <div class="donate">
    <h1>Support Our Project</h1>
    <p class="description">
      Thank you for using our dictionary application! If you find it useful, please consider supporting
      our project with a donation. Your contribution helps us maintain and improve the service.
    </p>

    <div class="donation-options">
      <div class="donation-card">
        <h3>One-time Donation</h3>
        <p>Support us with a single contribution of any amount.</p>
        <div class="amount-buttons">
          <button @click="setAmount(5)" :class="{ active: selectedAmount === 5 }">$5</button>
          <button @click="setAmount(10)" :class="{ active: selectedAmount === 10 }">$10</button>
          <button @click="setAmount(25)" :class="{ active: selectedAmount === 25 }">$25</button>
          <button @click="setAmount(50)" :class="{ active: selectedAmount === 50 }">$50</button>
        </div>
        <div class="custom-amount">
          <label for="custom-amount">Custom amount:</label>
          <input 
            type="number" 
            id="custom-amount" 
            v-model="customAmount" 
            @input="setCustomAmount"
            min="1"
            placeholder="Enter amount"
          />
        </div>
        <button class="donate-button" @click="processDonation">Donate Now</button>
      </div>

      <div class="donation-card">
        <h3>Monthly Support</h3>
        <p>Become a regular supporter with a monthly contribution.</p>
        <div class="amount-buttons">
          <button @click="setMonthlyAmount(3)" :class="{ active: selectedMonthlyAmount === 3 }">$3/mo</button>
          <button @click="setMonthlyAmount(5)" :class="{ active: selectedMonthlyAmount === 5 }">$5/mo</button>
          <button @click="setMonthlyAmount(10)" :class="{ active: selectedMonthlyAmount === 10 }">$10/mo</button>
          <button @click="setMonthlyAmount(20)" :class="{ active: selectedMonthlyAmount === 20 }">$20/mo</button>
        </div>
        <div class="custom-amount">
          <label for="custom-monthly-amount">Custom monthly amount:</label>
          <input 
            type="number" 
            id="custom-monthly-amount" 
            v-model="customMonthlyAmount" 
            @input="setCustomMonthlyAmount"
            min="1"
            placeholder="Enter amount"
          />
        </div>
        <button class="donate-button" @click="processMonthlyDonation">Subscribe</button>
      </div>
    </div>

    <div class="thank-you-section">
      <h2>Thank You to Our Supporters</h2>
      <p>We'd like to express our gratitude to everyone who has supported this project.</p>
      <div class="supporters">
        <div v-for="(supporter, index) in supporters" :key="index" class="supporter">
          {{ supporter.name }} - {{ supporter.tier }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from '@vue/runtime-dom'

const selectedAmount = ref<number | null>(null)
const customAmount = ref<number | null>(null)
const selectedMonthlyAmount = ref<number | null>(null)
const customMonthlyAmount = ref<number | null>(null)

// Real supporters would be fetched from an API
const supporters = ref([
  { name: 'Anonymous', tier: 'Gold Supporter' },
  { name: 'Anonymous', tier: 'Silver Supporter' },
  { name: 'Anonymous', tier: 'Bronze Supporter' }
])

const setAmount = (amount: number) => {
  selectedAmount.value = amount
  customAmount.value = null
}

const setCustomAmount = () => {
  if (customAmount.value) {
    selectedAmount.value = null
  }
}

const setMonthlyAmount = (amount: number) => {
  selectedMonthlyAmount.value = amount
  customMonthlyAmount.value = null
}

const setCustomMonthlyAmount = () => {
  if (customMonthlyAmount.value) {
    selectedMonthlyAmount.value = null
  }
}

const processDonation = () => {
  const amount = selectedAmount.value || customAmount.value
  if (!amount) {
    alert('Please select or enter a donation amount')
    return
  }
  
  // In a real application, this would connect to a payment processor
  alert(`Thank you for your one-time donation of $${amount}!`)
}

const processMonthlyDonation = () => {
  const amount = selectedMonthlyAmount.value || customMonthlyAmount.value
  if (!amount) {
    alert('Please select or enter a monthly donation amount')
    return
  }
  
  // In a real application, this would connect to a payment processor
  alert(`Thank you for your monthly donation of $${amount}!`)
}
</script>

<style scoped>
.donate {
  max-width: 900px;
  margin: 0 auto;
}

h1 {
  text-align: center;
  color: #2c3e50;
  margin-bottom: 20px;
}

.description {
  text-align: center;
  margin-bottom: 40px;
  font-size: 18px;
  line-height: 1.6;
  color: #555;
}

.donation-options {
  display: flex;
  justify-content: space-between;
  gap: 30px;
  margin-bottom: 50px;
}

.donation-card {
  flex: 1;
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 25px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.donation-card h3 {
  color: #2c3e50;
  margin-bottom: 15px;
  text-align: center;
}

.donation-card p {
  margin-bottom: 20px;
  text-align: center;
  color: #555;
}

.amount-buttons {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.amount-buttons button {
  flex: 1;
  margin: 0 5px;
  padding: 10px;
  background-color: #f1f1f1;
  border: 1px solid #ddd;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.amount-buttons button.active {
  background-color: #42b983;
  color: white;
  border-color: #42b983;
}

.custom-amount {
  margin-bottom: 20px;
}

.custom-amount label {
  display: block;
  margin-bottom: 5px;
  color: #555;
}

.custom-amount input {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
}

.donate-button {
  width: 100%;
  padding: 12px;
  background-color: #42b983;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.donate-button:hover {
  background-color: #3aa876;
}

.thank-you-section {
  margin-top: 50px;
  text-align: center;
}

.thank-you-section h2 {
  color: #2c3e50;
  margin-bottom: 15px;
}

.supporters {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 20px;
}

.supporter {
  background-color: #f1f1f1;
  padding: 10px 15px;
  margin: 5px;
  border-radius: 20px;
  font-size: 14px;
  color: #555;
}

@media (max-width: 768px) {
  .donation-options {
    flex-direction: column;
  }
}
</style> 