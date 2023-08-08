<script setup>
import { reactive, ref } from "vue";

const file = ref(null);
const upscaled = ref(null);

let upscaleHappened = ref(false);
let loading = ref(false);
let buttonLabel = ref("Upload a file");
let model = reactive({
  scale: 1,
  noise: 1,
  imageFile: null,
});

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
  fetch("https://scalar.hiyo.run/upscale", {
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
  downloadLink.download = `[${model.scale * 100}%][${model.noise}x]${
    model.imageFile[0].name
  }`;
  downloadLink.click();
}
</script>
<template>
  <div class="center">
    <input
      type="file"
      ref="file"
      id="upload"
      @change="readFile()"
      :disabled="loading"
    />
    <label :class="!loading ? 'button' : 'button-disabled'" for="upload">
      <span>{{ buttonLabel }}</span>
    </label>
    <span class="group">Scale</span>
    <div>
      <input
        type="radio"
        v-model="model.scale"
        value="1"
        :disabled="loading"
      />1
      <input
        type="radio"
        v-model="model.scale"
        value="2"
        :disabled="loading"
      />2
      <input
        type="radio"
        v-model="model.scale"
        value="4"
        :disabled="loading"
      />4
      <input
        type="radio"
        v-model="model.scale"
        value="8"
        :disabled="loading"
      />8
    </div>
    <span class="group">Noise Reduction</span>
    <div>
      <input
        type="radio"
        v-model="model.noise"
        value="1"
        :disabled="loading"
      />1
      <input
        type="radio"
        v-model="model.noise"
        value="2"
        :disabled="loading"
      />2
      <input
        type="radio"
        v-model="model.noise"
        value="3"
        :disabled="loading"
      />3
    </div>
    <div class="group">
      <button
        style="display: flex; align-items: center; justify-content: center"
        :class="!loading ? 'button' : 'button-disabled'"
        @click="submitUpscale"
      >
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
