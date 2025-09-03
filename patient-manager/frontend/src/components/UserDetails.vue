<template>
 <v-container class="d-flex align-center justify-center">
      <v-card width="800" rounded="lg" elevation="3" class="pa-8 ma-8" variant="tonal">
      <v-card-title class="text-h5 pb-3">User Details</v-card-title>

      <div v-if="user">
        <v-text-field v-model="user.firstName" density="compact" label="First Name" variant="underlined"
          prepend-inner-icon="mdi-account-outline" readonly class="mb-4"></v-text-field>

        <v-text-field v-model="user.lastName" density="compact" label="Last Name" variant="underlined"
          prepend-inner-icon="mdi-account-outline" readonly class="mb-4"></v-text-field>

        <v-text-field v-model="user.email" density="compact" label="Email" variant="underlined"
          prepend-inner-icon="mdi-email-outline" readonly class="mb-4"></v-text-field>

        <v-text-field v-model="user.role" density="compact" label="Role" variant="underlined"
          prepend-inner-icon="mdi-shield-account-outline" readonly class="mb-4"></v-text-field>
      </div>
      <div v-else>
        You don't have details
      </div>
    </v-card>
  </v-container>
</template>

<script lang="ts" setup>
import type { User } from '@/models/user';
import { useAuthStore } from '@/stores/auth';
import { onMounted } from 'vue';

const user = ref<User>()
const authStore = useAuthStore();

onMounted(async () => {
  const rez = authStore.User
  console.log(rez)
  user.value = rez
  console.log({ ...user.value })
  console.log(rez?.email)
});
</script>
