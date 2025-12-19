<script setup lang="ts"></script>

<template>
  <div class="app">
    <header class="header">
      <h1 class="title">codis</h1>
      <nav class="nav">
        <RouterLink to="/debug/discord-oauth" class="link">Discord OAuth Debug</RouterLink>
      </nav>
    </header>

    <main class="main">
      <RouterView />
    </main>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html,
body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
}

#app {
  width: 100%;
  height: 100%;
}
</style>

<style scoped>
.app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #020617;
  color: #e5e7eb;
  font-family:
    system-ui,
    -apple-system,
    BlinkMacSystemFont,
    'SF Pro Text',
    sans-serif;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1.5rem;
  border-bottom: 1px solid rgba(31, 41, 55, 0.9);
  background: rgba(15, 23, 42, 0.96);
}

.title {
  font-size: 1.1rem;
  font-weight: 600;
}

.nav {
  display: flex;
  gap: 0.75rem;
}

.link {
  font-size: 0.9rem;
  color: #9ca3af;
  text-decoration: none;
}

.link:hover {
  color: #e5e7eb;
}

.main {
  flex: 1;
}
</style>
