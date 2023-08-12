<script setup>
import { Menu, MenuButton, MenuItem, MenuItems } from '@headlessui/vue'
import { Bars3Icon } from '@heroicons/vue/24/outline'
import { onMounted } from 'vue';
import { RouterView, useRouter } from 'vue-router'
const router = useRouter()

function changePage(route) {
  router.push(route)
}
</script>

<template>
  <Menu as="div" class="absolute top-4 right-4">
    <div>
      <MenuButton
        class="relative flex rounded-full bg-gray-800 text-sm focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-gray-800">
        <span class="absolute -inset-1.5" />
        <span class="sr-only">Open user menu</span>
        <Bars3Icon class="h-8 w-8 rounded-full text-arisu-800 dark:text-arisu-200" />
      </MenuButton>
    </div>
    <transition enter-active-class="transition ease-out duration-100" enter-from-class="transform opacity-0 scale-95"
      enter-to-class="transform opacity-100 scale-100" leave-active-class="transition ease-in duration-75"
      leave-from-class="transform opacity-100 scale-100" leave-to-class="transform opacity-0 scale-95">
      <MenuItems
        class="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
        <MenuItem v-slot="{ active }" @click="changePage('/')">
        <a href="#"
          :class="[active ? 'bg-gray-100' : '', 'block px-4 py-2 text-sm text-gray-700 hover:bg-arisu-500']">Upload</a>
        </MenuItem>
        <MenuItem v-slot="{ active }" @click="changePage('/images')">
        <a href="#"
          :class="[active ? 'bg-gray-100' : '', 'block px-4 py-2 text-sm text-gray-700 hover:bg-arisu-500']">Your
          Images</a>
        </MenuItem>
        <MenuItem v-slot="{ active }" @click="changePage('/session-id')">
        <a href="#"
          :class="[active ? 'bg-gray-100' : '', 'block px-4 py-2 text-sm text-gray-700 hover:bg-arisu-500']">Session
          ID</a>
        </MenuItem>
      </MenuItems>
    </transition>
  </Menu>
  <RouterView /></template>

<style scoped></style>
