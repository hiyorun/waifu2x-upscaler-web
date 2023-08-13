<script setup>
import { reactive, ref, onMounted, onBeforeUnmount } from "vue";
import { useWebSocket } from "@/composables/useWebSocket";
import { useSession } from "@/states/sessionid";
import { useImageHelper } from "@/composables/useImageHelper";
import { useApiUrl } from '../composables/useAPI';

const { socket, createWebSocket, handleConnection, handleDisconnection, handleError, sendHeartbeat } = useWebSocket();
const { downloadImage } = useImageHelper()
const sessionStore = useSession()

let images = reactive({ entries: [] })

const heartbeatInterval = setInterval(sendHeartbeat, 5000);

onMounted(
  () => {
    sessionStore.getSession
    createWebSocket(getImages())

    socket.value.onmessage = (event) => {
      console.log(event)
      if (event.data === sessionStore.getSession) {
        getImages()
      }
    };
  }
)

onBeforeUnmount(() => {
  clearInterval(heartbeatInterval);
  if (socket.value) {
    socket.value.removeEventListener('open', handleConnection);
    socket.value.removeEventListener('close', handleDisconnection);
    socket.value.removeEventListener('error', handleError);
    socket.value.close();
  }
});

async function getImages() {
  const response = await fetch(useApiUrl('/get-images?') +
    new URLSearchParams({
      uuid: sessionStore.getSession,
    })
  );
  const jsonData = await response.json();
  images.entries = jsonData;
}

</script>
<template>
  <div class="flex flex-1 flex-col justify-center px-6 py-12 lg:px-8 bg-arisu-200 dark:bg-arisu-900 noisy min-h-screen min-w-screen">
    <p class="text-sm font-semibold text-arisu-900 dark:text-arisu-100">Click on the list to download</p>
    <ul role="list">
      <li v-for="entry in images.entries" :key="entry.uuid"
        class="rounded-full flex justify-between gap-x-6 py-5 my-2 px-5 hover:bg-arisu-300 dark:hover:bg-arisu-700"
        @click="downloadImage(entry.filename, entry.status)">
        <div class="flex min-w-0 gap-x-4">
          <!-- <img class="h-12 w-12 flex-none rounded-full bg-arisu-50" :src="entry.imageUrl" alt="" /> -->
          <div class="min-w-0 flex-auto">
            <p class="text-sm font-semibold text-arisu-900 dark:text-arisu-100">{{ entry.name }}</p>
          </div>
        </div>
        <div class="hidden shrink-0 sm:flex sm:flex-col sm:items-end">
          <p v-if="entry.status !== 'done'" class="mt-1 text-xs leading-5 text-arisu-500">
            {{ entry.status }}
          </p>
          <div v-else class="mt-1 flex items-center gap-x-1.5">
            <div class="flex-none rounded-full bg-arisu-600/20 dark:bg-arisu-300/20 p-1">
              <div class="h-1.5 w-1.5 rounded-full bg-arisu-600 dark:bg-arisu-300" />
            </div>
            <p class="text-xs leading-5 text-arisu-800 dark:text-arisu-200">{{ entry.status }}</p>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>