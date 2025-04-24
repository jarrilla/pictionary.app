<template>
  <div class="tip">
    <h1>Support Our Project</h1>
    <p class="description">
      Thank you for using our dictionary application! If you find it useful, please consider supporting
      our project with a tip. Your contribution helps us maintain and improve the service.
    </p>

    <div class="donation-button-container">
      <a 
        href="https://buy.stripe.com/00g2bi0WR7Q5448dQQ" 
        target="_blank" 
        rel="noopener noreferrer"
        class="tip-button"
      >
        Tip the Project
      </a>
    </div>

    <div class="thank-you-section">
      <h2>Thank you for your support!</h2>
    </div>

    <div class="supporters">
      <p>This project is dedicated to SEM, my creative muse and best friend.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.supporters {
  position: absolute;
  left: calc(50% - 14rem);
  bottom: 1rem;
  text-align: center;
  width: 28rem;
  font-size: 0.8rem;
  font-family: 'Birthstone', cursive;
  font-style: italic;
}

.tip {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
  text-align: center;
}

h1 {
  color: #2c3e50;
  margin-bottom: 20px;
}

.description {
  margin-bottom: 40px;
  font-size: 18px;
  line-height: 1.6;
  color: #555;
}

.donation-button-container {
  margin: 40px 0;
}

.tip-button {
  display: inline-block;
  padding: 15px 40px;
  background-color: #42b983;
  color: white;
  text-decoration: none;
  border-radius: 4px;
  font-size: 18px;
  transition: background-color 0.3s;
}

.tip-button:hover {
  background-color: #3aa876;
}

.thank-you-section {
  margin-top: 50px;
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
</style> 