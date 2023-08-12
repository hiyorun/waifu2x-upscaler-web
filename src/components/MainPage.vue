<script setup>
import { reactive, ref, onMounted, onBeforeUnmount } from "vue";
import { useWebSocket } from "@/composables/useWebSocket";
import { useImageHelper } from "@/composables/useImageHelper";
import { useBusy } from '@/states/busy.js';

const { socket, createWebSocket, handleConnection, handleDisconnection, handleError, sendHeartbeat } = useWebSocket();
const { uploadImage, downloadImage, getImages, } = useImageHelper

const busy = useBusy()
const file = ref(null);
let session = reactive({ uuid: "" })
let sIDExist = ref(false)
let modify = ref(false)
let buttonLabel = ref("Upload a file");
let images = reactive({ entries: [] })
let model = reactive({
  scale: 1,
  noise: 1,
  imageFile: null,
});
const heartbeatInterval = setInterval(sendHeartbeat, 5000);

onMounted(
  () => {
    session.uuid = localStorage.getItem('session_id')
    if (session.uuid) {
      sIDExist.value = true
    }

    createWebSocket()

    socket.value.onmessage = (event) => {
      console.log(event)
      if (event.data === session.uuid) {
        getImages();
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

function generateUUID() {
  session.uuid = ""
  handleUUID()
}

function handleUUID() {
  console.log(session.uuid)
  if (!session.uuid) {
    session.uuid = crypto.randomUUID()
  }
  localStorage.setItem('session_id', session.uuid);
  sIDExist.value = true
  modify.value = false
}

function cancelModify() {
  modify.value = false
  session.uuid = localStorage.getItem('session_id')
}

function readFile() {
  model.imageFile = file.value.files;
  buttonLabel.value = "Current = " + model.imageFile[0].name;
}

function upload() {
  busy.setBusy(true)
  uploadImage(model)
    .catch((err) => {
      console.error(err);
      return err
    })
    .finally(()=>busy.setBusy(false))
}
</script>
<template>
  <div v-if="!sIDExist" class="dialogbackground">
    <div class="dialog">
      <span>Input your existing session UUID or create a new one</span>
      <div>
        <input type="text" v-model="session.uuid" placeholder="Session ID">
        <button class="button" @click="handleUUID()">Save</button>
      </div>
      <button class="button" @click="generateUUID()">Generate new session UUID</button>
    </div>
  </div>
  <div v-if="modify" class="dialogbackground">
    <div class="dialog">
      <span>Modify or regenerate your session UUID</span>
      <div>
        <input type="text" v-model="session.uuid" placeholder="Session ID">
        <button class="button" @click="handleUUID()">Save</button>
      </div>
      <button class="button" @click="generateUUID()">Generate new session UUID</button>
      <button class="button" @click="cancelModify()">Cancel</button>
    </div>
  </div>
  <div class="center">
    <span>Your session UUID is: "{{ session.uuid }}"</span>
    <div style="margin: 1em 0;">
      <button @click="() => { modify = true }" class="button">Configure</button>
    </div>
    <input type="file" ref="file" id="upload" @change="readFile()" :disabled="busy.isBusy" />
    <label :class="!busy.isBusy ? 'button' : 'button-disabled'" for="upload">
      <span>{{ buttonLabel }}</span>
    </label>
    <span class="group">Scale</span>
    <div>
      <input type="radio" v-model="model.scale" value="1" :disabled="busy.isBusy" />1
      <input type="radio" v-model="model.scale" value="2" :disabled="busy.isBusy" />2
      <input type="radio" v-model="model.scale" value="4" :disabled="busy.isBusy" />4
      <input type="radio" v-model="model.scale" value="8" :disabled="busy.isBusy" />8
    </div>
    <span class="group">Noise Reduction</span>
    <div>
      <input type="radio" v-model="model.noise" value="1" :disabled="busy.isBusy" />1
      <input type="radio" v-model="model.noise" value="2" :disabled="busy.isBusy" />2
      <input type="radio" v-model="model.noise" value="3" :disabled="busy.isBusy" />3
    </div>
    <div class="group">
      <button style="display: flex; align-items: center; justify-content: center"
        :class="!busy.isBusy ? 'button' : 'button-disabled'" @click="upload()">
        <span v-if="busy.isBusy"> Uploading </span>
        <span v-else>Post</span>
        <div v-if="busy.isBusy" class="lds-ripple">
          <div></div>
          <div></div>
        </div>
      </button>
    </div>
    <span class="group">Your Images</span>
    <table class="group table">
      <thead>
        <tr>
          <td>Name</td>
          <td>Status</td>
          <td colspan="2" width="30%">Action</td>
        </tr>
      </thead>
      <tbody>
        <tr v-for="entry in images.entries" :key="entry.uuid">
          <td>{{ entry.name }}</td>
          <td>{{ entry.status }}</td>
          <td>
            <button class="button" @click="downloadImage(entry.filename)">
              <span class="material-symbols-outlined">
                download
              </span>
            </button>
          </td>
          <td>
            <button class="button" @click="downloadImage(entry.filename)">
              <span class="material-symbols-outlined">
                delete
              </span>
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
<style scoped>
.table {
  width: 70%;
}

.table td {
  padding: 0.5em;
  text-align: center;
}

table>thead {
  background-color: #ff9ead;
}

.table>tbody>tr:nth-child(even) {
  background-color: #ff9ead;
}

.dialogbackground {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: #fff8;
}

.dialog {
  display: flex;
  flex-direction: column;
  position: absolute;
  top: 50%;
  left: 50%;
  background-color: #ffdbde;
  padding: .5em 2em 2em 2em;
  transform: translate(-50%, -50%);
}

.dialog>* {
  margin-top: 1.5em;
}

.imgContainer {
  position: relative;
  overflow: hidden;
}

.download {
  width: 100%;
  height: 400px;
  position: absolute;
  top: 0;
  left: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  background-color: #fff7;
  cursor: pointer;
}

.download:hover {
  opacity: 1;
}

.upscaled {
  width: auto;
  height: 400px;
}
</style>
