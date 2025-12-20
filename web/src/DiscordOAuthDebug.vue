<script setup lang="ts">
import { computed, ref } from 'vue'

const BACKEND_BASE_URL = 'http://localhost:80'

const inviteLinkResponse = ref<unknown | null>(null)
const inviteLinkError = ref<string | null>(null)

const queryParams = computed(() => {
  const params = new URLSearchParams(window.location.search)
  const entries: { key: string; value: string }[] = []
  params.forEach((value, key) => {
    entries.push({ key, value })
  })
  return entries
})

function initialCallbackBody(): string {
  const params = new URLSearchParams(window.location.search)
  const code = params.get('code')
  const state = params.get('state')

  if (code || state) {
    return JSON.stringify(
      {
        code: code ?? 'PASTE_CODE_HERE',
        state: state ?? 'PASTE_STATE_HERE',
      },
      null,
      2,
    )
  }

  return '{"code": "PASTE_CODE_HERE", "state": "PASTE_STATE_HERE"}'
}

const callbackBody = ref<string>(initialCallbackBody())
const callbackResponse = ref<unknown | null>(null)
const callbackError = ref<string | null>(null)

const helloworldResponse = ref<unknown | null>(null)
const helloworldError = ref<string | null>(null)

const inviteLinkUrl = `${BACKEND_BASE_URL}/discord/invite_link`
const callbackUrl = `${BACKEND_BASE_URL}/discord/callback`
const helloworldUrl = `${BACKEND_BASE_URL}/discord/guilds`

async function fetchInviteLink() {
  inviteLinkError.value = null
  inviteLinkResponse.value = null
  try {
    const res = await fetch(inviteLinkUrl, {
      credentials: 'include',
    })
    const json = await res.json()
    inviteLinkResponse.value = json
    console.log('Invite link response', json)
  } catch (err: unknown) {
    if (err instanceof Error) {
      const message = err?.message ?? String(err)
      inviteLinkError.value = message
      console.error('Error fetching invite link', err)
    }
  }
}

async function reset() {
  window.location.href = window.location.pathname
}

async function redirectToInviteLink() {
  inviteLinkError.value = null
  try {
    const res = await fetch(inviteLinkUrl, {
      credentials: 'include',
    })
    const json = await res.json()
    const inviteLink = (json as { discord_invite_link?: string }).discord_invite_link
    if (inviteLink) {
      window.location.href = inviteLink
    } else {
      inviteLinkError.value = 'No discord_invite_link found in response'
      console.error('Invalid response format', json)
    }
  } catch (err: unknown) {
    if (err instanceof Error) {
      const message = err?.message ?? String(err)
      inviteLinkError.value = message
      console.error('Error fetching invite link', err)
    }
  }
}

async function postCallbackBody() {
  callbackError.value = null
  callbackResponse.value = null
  try {
    const res = await fetch(callbackUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: 'include',
      body: callbackBody.value,
    })
    const text = await res.text()
    let parsed: unknown = text
    try {
      parsed = JSON.parse(text)
    } catch {
      // keep plain text
    }
    callbackResponse.value = parsed
    console.log('Callback response', parsed)
  } catch (err: unknown) {
    if (err instanceof Error) {
      const message = err?.message ?? String(err)
      callbackError.value = message
      console.error('Error posting callback body', err)
    }
  }
}

async function fetchHelloworld() {
  helloworldError.value = null
  helloworldResponse.value = null
  try {
    const res = await fetch(helloworldUrl, {
      credentials: 'include',
    })
    const text = await res.text()
    let parsed: unknown = text
    try {
      parsed = JSON.parse(text)
    } catch {
      // keep plain text
    }
    helloworldResponse.value = parsed
    console.log('Helloworld response', parsed)
  } catch (err: unknown) {
    if (err instanceof Error) {
      const message = err?.message ?? String(err)
      helloworldError.value = message
      console.error('Error fetching helloworld', err)
    }
  }
}
</script>

<template>
  <main class="page">
    <section class="card">
      <h1>Discord OAuth Debug</h1>
      <p class="subtitle">
        Backend base:
        <code>{{ BACKEND_BASE_URL }}</code>
      </p>

      <div style="display: flex; gap: 0.75rem; margin-bottom: 0.75rem">
        <button type="button" class="btn" @click="reset">Reset</button>
      </div>

      <section class="block">
        <h2>1. Query params on this page</h2>
        <p class="hint">
          If your backend redirects back to this frontend, copy the full URL and open it here to
          inspect
          <code>code</code>, <code>state</code>, or <code>error</code> parameters.
        </p>
        <table v-if="queryParams.length" class="table">
          <thead>
            <tr>
              <th>Key</th>
              <th>Value</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in queryParams" :key="item.key">
              <td>{{ item.key }}</td>
              <td>{{ item.value }}</td>
            </tr>
          </tbody>
        </table>
        <p v-else class="muted">No query parameters detected on this URL.</p>
      </section>

      <section class="block">
        <h2>2. Get Discord invite link</h2>
        <p class="hint">
          Calls
          <code>GET /discord/invite_link</code>
          and shows the JSON response.
        </p>
        <div style="display: flex; gap: 0.75rem; margin-bottom: 0.75rem">
          <button type="button" class="btn" @click="fetchInviteLink">
            Call /discord/invite_link
          </button>
          <button type="button" class="btn" @click="redirectToInviteLink">
            Redirect to invite link
          </button>
        </div>
        <div class="result">
          <p class="label">Result</p>
          <pre v-if="inviteLinkResponse">{{ JSON.stringify(inviteLinkResponse, null, 2) }}</pre>
          <p v-else class="muted">No response yet.</p>
          <p v-if="inviteLinkError" class="error">Error: {{ inviteLinkError }}</p>
        </div>
      </section>

      <section class="block">
        <h2>3. POST to /discord/callback</h2>
        <p class="hint">
          Sends a JSON body to
          <code>POST /discord/callback</code>
          . Adjust the
          <code>code</code>
          and
          <code>state</code>
          below.
        </p>
        <textarea v-model="callbackBody" class="textarea" rows="6" />
        <button type="button" class="btn" @click="postCallbackBody">Send callback body</button>
        <div class="result">
          <p class="label">Result</p>
          <pre v-if="callbackResponse">{{ JSON.stringify(callbackResponse, null, 2) }}</pre>
          <p v-else class="muted">No response yet.</p>
          <p v-if="callbackError" class="error">Error: {{ callbackError }}</p>
        </div>
      </section>

      <section class="block">
        <h2>4. GET /helloworld</h2>
        <p class="hint">
          Calls
          <code>GET /helloworld</code>
          and shows the response.
        </p>
        <button type="button" class="btn" @click="fetchHelloworld">Call /helloworld</button>
        <div class="result">
          <p class="label">Result</p>
          <pre v-if="helloworldResponse !== null">{{
            JSON.stringify(helloworldResponse, null, 2)
          }}</pre>
          <p v-else class="muted">No response yet.</p>
          <p v-if="helloworldError" class="error">Error: {{ helloworldError }}</p>
        </div>
      </section>
    </section>
  </main>
</template>

<style scoped>
.page {
  min-height: 100vh;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding: 3rem 1.5rem;
  background: radial-gradient(circle at top left, #1f2933, #020617);
  color: #e5e7eb;
  font-family:
    system-ui,
    -apple-system,
    BlinkMacSystemFont,
    'SF Pro Text',
    sans-serif;
}

.card {
  width: 100%;
  max-width: 960px;
  background: rgba(15, 23, 42, 0.96);
  border-radius: 1rem;
  padding: 2rem 2.25rem;
  box-shadow:
    0 24px 60px rgba(0, 0, 0, 0.6),
    0 0 0 1px rgba(148, 163, 184, 0.2);
}

h1 {
  font-size: 1.75rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.subtitle {
  margin-bottom: 1.75rem;
  color: #9ca3af;
  font-size: 0.9rem;
}

.block {
  padding-top: 1.5rem;
  margin-top: 1.5rem;
  border-top: 1px solid rgba(55, 65, 81, 0.9);
}

.block:first-of-type {
  border-top: none;
  padding-top: 0;
  margin-top: 0;
}

h2 {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 0.35rem;
}

.hint {
  font-size: 0.9rem;
  color: #9ca3af;
  margin-bottom: 0.75rem;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.35rem;
  padding: 0.55rem 1.1rem;
  border-radius: 999px;
  border: none;
  cursor: pointer;
  background: linear-gradient(135deg, #5865f2, #3b82f6);
  color: white;
  font-size: 0.9rem;
  font-weight: 500;
  box-shadow:
    0 10px 25px rgba(59, 130, 246, 0.45),
    0 0 0 1px rgba(59, 130, 246, 0.6);
  transition:
    transform 0.08s ease-out,
    box-shadow 0.08s ease-out,
    filter 0.12s ease-out;
}

.btn:hover {
  transform: translateY(-1px);
  filter: brightness(1.05);
  box-shadow:
    0 14px 35px rgba(59, 130, 246, 0.6),
    0 0 0 1px rgba(96, 165, 250, 0.9);
}

.btn:active {
  transform: translateY(0);
  box-shadow:
    0 6px 18px rgba(37, 99, 235, 0.6),
    0 0 0 1px rgba(59, 130, 246, 0.6);
}

.textarea {
  width: 100%;
  margin-bottom: 0.75rem;
  border-radius: 0.75rem;
  border: 1px solid rgba(55, 65, 81, 0.95);
  background: rgba(15, 23, 42, 0.85);
  color: #e5e7eb;
  padding: 0.6rem 0.75rem;
  font-family:
    ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New',
    monospace;
  font-size: 0.85rem;
  resize: vertical;
}

.textarea:focus {
  outline: none;
  border-color: rgba(129, 140, 248, 0.9);
  box-shadow: 0 0 0 1px rgba(129, 140, 248, 0.9);
}

.result {
  margin-top: 0.75rem;
}

.label {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: #9ca3af;
  margin-bottom: 0.25rem;
}

pre {
  background: rgba(15, 23, 42, 0.9);
  border-radius: 0.75rem;
  padding: 0.75rem 0.9rem;
  font-size: 0.8rem;
  overflow-x: auto;
  border: 1px solid rgba(55, 65, 81, 0.95);
}

.muted {
  font-size: 0.85rem;
  color: #6b7280;
}

.error {
  margin-top: 0.35rem;
  font-size: 0.85rem;
  color: #f97373;
}

.table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 0.25rem;
  font-size: 0.85rem;
}

.table th,
.table td {
  border: 1px solid rgba(55, 65, 81, 0.9);
  padding: 0.3rem 0.5rem;
}

.table th {
  background: rgba(15, 23, 42, 0.9);
  text-align: left;
  font-weight: 500;
}
</style>
