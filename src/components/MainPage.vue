<script setup>
import { reactive, ref, onMounted } from "vue";

const file = ref(null);
const upscaled = ref(null);

let session = reactive({ uuid: "" })
let sIDExist = ref(false)
let modify = ref(false)
let upscaleHappened = ref(false);
let loading = ref(false);
let buttonLabel = ref("Upload a file");
let model = reactive({
  scale: 1,
  noise: 1,
  imageFile: null,
});

onMounted(
  () => {
    session.uuid = localStorage.getItem('session_id')
    if (session.uuid) {
      sIDExist.value = true
    }
    console.log(session.uuid)
  }
)

function generateID() {
  session.uuid = crypto.randomUUID()
  console.log(session)
  localStorage.setItem('session_id', session.uuid);
  sIDExist.value = true
  modify.value = false
}

function storeSessionID() {
  console.log(session.uuid)
  if (!session.uuid) {
    return
  }
  localStorage.setItem('session_id', session);
  modify.value = false
}

function configureSessionID() {
  modify.value = true
}

function cancelModify() {
  modify.value = false
  session.uuid = localStorage.getItem('session_id')
}

function readFile() {
  model.imageFile = file.value.files;
  buttonLabel.value = "Current = " + model.imageFile[0].name;
}

function submitUpscale() {
  upscaleHappened.value = false;
  loading.value = true;
  let form = new FormData();
  form.append("imageFile", model.imageFile[0]);
  form.append("scale", model.scale);
  form.append("noise", model.noise);
  form.append("uuid",session.uuid)
  fetch("http://127.0.0.1:8080/upload", {
    method: "POST",
    body: form,
  })
    .then(async (resp) => {
      let file = await resp.blob();
      let reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onloadend = () => {
        // console.log(reader.result);
        upscaled.value.src = reader.result;
      };
    })
    .catch((err) => {
      console.error(err);
    })
    .finally(() => {
      loading.value = false;
      upscaleHappened.value = true;
    });
}

function download() {
  // console.log(upscaled.value, fileName);
  let linkSource = upscaled.value.src;
  let downloadLink = document.createElement("a");
  downloadLink.href = linkSource;
  downloadLink.download = `[${model.scale * 100}%][${model.noise}x]${model.imageFile[0].name
    }`;
  downloadLink.click();
}
</script>
<template>
  <div v-if="!sIDExist" class="dialogbackground">
    <div class="dialog">
      <span>Input your existing session UUID or create a new one</span>
      <div>
        <input type="text" v-model="session.uuid" placeholder="Session ID">
        <button class="button" @click="storeSessionID()">Save</button>
      </div>
      <button class="button" @click="generateID()">Generate new session UUID</button>
    </div>
  </div>
  <div v-if="modify" class="dialogbackground">
    <div class="dialog">
      <span>Modify or regenerate your session UUID</span>
      <div>
        <input type="text" v-model="session.uuid" placeholder="Session ID">
        <button class="button" @click="storeSessionID()">Save</button>
      </div>
      <button class="button" @click="generateID()">Generate new session UUID</button>
      <button class="button" @click="cancelModify()">Cancel</button>
    </div>
  </div>
  <div class="center">
    <span>Your session UUID is: "{{ session.uuid }}"</span>
    <div style="margin: 1em 0;">
      <button @click="configureSessionID()" class="button">Configure</button>
    </div>
    <input type="file" ref="file" id="upload" @change="readFile()" :disabled="loading" />
    <label :class="!loading ? 'button' : 'button-disabled'" for="upload">
      <span>{{ buttonLabel }}</span>
    </label>
    <span class="group">Scale</span>
    <div>
      <input type="radio" v-model="model.scale" value="1" :disabled="loading" />1
      <input type="radio" v-model="model.scale" value="2" :disabled="loading" />2
      <input type="radio" v-model="model.scale" value="4" :disabled="loading" />4
      <input type="radio" v-model="model.scale" value="8" :disabled="loading" />8
    </div>
    <span class="group">Noise Reduction</span>
    <div>
      <input type="radio" v-model="model.noise" value="1" :disabled="loading" />1
      <input type="radio" v-model="model.noise" value="2" :disabled="loading" />2
      <input type="radio" v-model="model.noise" value="3" :disabled="loading" />3
    </div>
    <div class="group">
      <button style="display: flex; align-items: center; justify-content: center"
        :class="!loading ? 'button' : 'button-disabled'" @click="submitUpscale">
        <span v-if="loading"> Upscaling </span>
        <span v-else>Post</span>
        <div v-if="loading" class="lds-ripple">
          <div></div>
          <div></div>
        </div>
      </button>
    </div>
    <div v-if="upscaleHappened" class="imgContainer group">
      <span class="download" @click="download">
        <span>Click to download</span>
      </span>
      <img ref="upscaled" class="upscaled" />
    </div>
  </div>
</template>
<style scoped>
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
