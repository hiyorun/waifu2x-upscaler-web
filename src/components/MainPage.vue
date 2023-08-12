<script setup>
import { reactive, ref } from "vue";
import { useImageHelper } from "@/composables/useImageHelper";
import { useSession } from "@/states/sessionid";

import { useBusy } from '@/states/busy.js';

const { uploadImage } = useImageHelper()
const sessionStore = useSession()

const busy = useBusy()
const file = ref(null);
let model = reactive({
  scale: 1,
  noise: 1,
  imageFile: null,
  uuid: ""
});

function readFile() {
  model.imageFile = file.value.files;
}

function upload() {
  busy.setBusy(true)
  model.uuid = sessionStore.getSession
  if(model.uuid === ""){
    console.log("No UUID")
    return
  }
  uploadImage(model)
    .catch((err) => {
      console.error(err);
      return err
    })
    .then(() => {
      model = reactive({
        scale: 1,
        noise: 1,
        imageFile: null,
        uuid: sessionStore.getSession
      });
    })
    .finally(() => {
      busy.setBusy(false)
    })
}
</script>
<template>
  <div class="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8 ">
    <form>
      <div class="space-y-12">
        <h2 class="text-base font-semibold leading-7 text-arisu-900 dark:text-arisu-100">Image to process</h2>
        <p class="mt-1 text-sm text-arisu-700 dark:text-arisu-200">Pick an image and select your desired scale and
          denoising
          ratio.
          Image won't be processed until you press upload.</p>
        <div class="mt-10 grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
          <div class="sm:col-span-4">
            <label for="imagefile" class="block text-sm font-medium text-arisu-900 dark:text-arisu-100">Select an
              image</label>
            <div class="mt-2">
              <div
                class="flex rounded-md shadow-sm ring-1 ring-inset ring-arisu-300 focus-within:ring-2 focus-within:ring-inset focus-within:ring-arisu-600 sm:max-w-md">
                <input type="file" ref="file" name="imagefile" id="imagefile" @change="readFile()"
                  class="block flex-1 border-0 bg-transparent py-1.5 pl-1 text-arisu-900 dark:text-arisu-100 placeholder:text-arisu-400 focus:ring-0 sm:text-sm sm"
                  :disabled="busy.isBusy" />
              </div>
            </div>
          </div>
        </div>

        <div class="mt-10 space-y-10">
          <fieldset>
            <legend class="text-sm font-semibold text-arisu-900 dark:text-arisu-100">Scale</legend>
            <p class="mt-1 text-sm text-arisu-700 dark:text-arisu-200">Choose the percentage by which your image will
              be
              enlarged.
            </p>
            <div class="mt-6 space-y-3">
              <div class="flex items-center gap-x-3">
                <input id="1x" value="1" v-model="model.scale" type="radio"
                  class="h-4 w-4 border-arisu-300 text-arisu-700 dark:text-arisu-200 focus:ring-arisu-600" />
                <label for="1x" class="block text-sm font-medium text-arisu-900 dark:text-arisu-100">Keep</label>
              </div>
              <div class="flex items-center gap-x-3">
                <input id="2x" value="2" v-model="model.scale" type="radio"
                  class="h-4 w-4 border-arisu-300 text-arisu-700 dark:text-arisu-200 focus:ring-arisu-600" />
                <label for="2x" class="block text-sm font-medium text-arisu-900 dark:text-arisu-100">200%</label>
              </div>
              <div class="flex items-center gap-x-3">
                <input id="4x" value="4" v-model="model.scale" type="radio"
                  class="h-4 w-4 border-arisu-300 text-arisu-700 dark:text-arisu-200 focus:ring-arisu-600" />
                <label for="4x" class="block text-sm font-medium text-arisu-900 dark:text-arisu-100">400%</label>
              </div>
              <div class="flex items-center gap-x-3">
                <input id="8x" value="8" v-model="model.scale" type="radio"
                  class="h-4 w-4 border-arisu-300 text-arisu-700 dark:text-arisu-200 focus:ring-arisu-600" />
                <label for="8x" class="block text-sm font-medium text-arisu-900 dark:text-arisu-100">800%</label>
              </div>
            </div>
          </fieldset>
          <fieldset>
            <legend class="text-sm font-semibold text-arisu-900 dark:text-arisu-100">Denoise</legend>
            <p class="mt-1 text-sm text-arisu-700 dark:text-arisu-200">How strong do you want your image to be
              denoised.</p>
            <div class="mt-6 space-y-3">
              <div class="flex items-center gap-x-3">
                <input id="noiseweak" value="1" v-model="model.noise" type="radio"
                  class="h-4 w-4 border-arisu-300 text-arisu-700 dark:text-arisu-200 focus:ring-arisu-600" />
                <label for="noiseweak" class="block text-sm font-medium text-arisu-900 dark:text-arisu-100">Weak</label>
              </div>
              <div class="flex items-center gap-x-3">
                <input id="noisefutsuu" value="2" v-model="model.noise" type="radio"
                  class="h-4 w-4 border-arisu-300 text-arisu-700 dark:text-arisu-200 focus:ring-arisu-600" />
                <label for="noisefutsuu"
                  class="block text-sm font-medium text-arisu-900 dark:text-arisu-100">Normal</label>
              </div>
              <div class="flex items-center gap-x-3">
                <input id="noisestrong" value="3" v-model="model.noise" type="radio"
                  class="h-4 w-4 border-arisu-300 text-arisu-700 dark:text-arisu-200 focus:ring-arisu-600" />
                <label for="noisestrong"
                  class="block text-sm font-medium text-arisu-900 dark:text-arisu-100">Stronk</label>
              </div>
            </div>
          </fieldset>
        </div>
      </div>
      <div class="mt-6 flex items-center justify-end gap-x-6">
        <button type="button" @click="upload()"
          class="rounded-md bg-arisu-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-arisu-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-arisu-600 "
          :disabled="busy.isBusy">Upload</button>
      </div>
    </form>
  </div>
</template>