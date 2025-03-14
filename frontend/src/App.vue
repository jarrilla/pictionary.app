<script setup lang="ts">
</script>

<template>
  <header>
    <nav>
      <router-link to="/">Home</router-link> |
      <router-link to="/donate">Donate</router-link>
    </nav>
  </header>
  <main>
    <router-view />
  </main>
</template>

<style>
#app {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

header {
  margin-bottom: 20px;
}

nav {
  padding: 15px 0;
  text-align: center;
}

nav a {
  font-weight: bold;
  color: var(--secondary-color);
  text-decoration: none;
  margin: 0 10px;
}

nav a.router-link-exact-active {
  color: var(--primary-color);
}

main {
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}
</style>
