<script setup>
import { useSession } from '@/states/sessionid';
import { onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router'

const session = reactive({ uuid: "" })
const router = useRouter()

const sessionStore = useSession()

onMounted(() => {
  session.uuid = sessionStore.getSession
  console.log(session.uuid)
})

function generateUUID() {
  sessionStore.setSession(crypto.randomUUID())
  if (sessionStore.getSession !== "") {
    router.push('/')
  }
}

function updateUUID() {
  sessionStore.setSession(session.uuid)
  console.log("Update UUID")
  if (sessionStore.getSession !== "") {
    router.push('/')
  }
}

</script>
<template>
  <div class="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8 ">
    <form>
      <div>
        <h2 class="text-base font-semibold leading-7 text-arisu-900 dark:text-arisu-100">Session ID</h2>
        <p class="mt-1 text-sm text-arisu-700 dark:text-arisu-200">Please generate a new session ID, or enter your old
          one. This will be only used by you to track image processing progress. No identity other than the session ID is
          collected.</p>
        <div class="mt-10 grid grid-cols-1 sm:grid-cols-6">
          <div class="sm:col-span-4">
            <div
              class="flex rounded-md shadow-sm ring-1 ring-inset ring-arisu-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-arisu-600 sm:max-w-md">
              <input type="text" name="sessionid" id="sessionid" autocomplete="sessionid" v-model="session.uuid"
                class="block flex-1 border-0 bg-transparent py-1.5 pl-1 text-arisu-900 dark:text-arisu-100 placeholder:text-gray-400 focus:ring-0 sm:text-sm sm:leading-6"
                placeholder="Session ID" />
            </div>
          </div>
        </div>
      </div>
      <div class="mt-6 flex items-center justify-end gap-x-6">
        <button type="button" class="text-sm font-semibold leading-6 text-arisu-800 dark:text-arisu-200 shadow-sm"
          @click="generateUUID()">Generate New
          Session ID</button>
        <button type="button" @click="updateUUID()" :disabled="session.uuid === ''"
          class="rounded-md bg-arisu-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-arisu-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-arisu-600">Save</button>
      </div>
    </form>
  </div>
</template>