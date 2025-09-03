<template>
    <v-container>
        <v-btn @click="goBack" color="grey" class="mb-4">
            <v-icon start>mdi-arrow-left</v-icon>
            Back to list
        </v-btn>

        <v-card v-if="patient" :loading="isLoading">
            <v-card-title class="pa-4">
                <span class="text-h5">Patient Details</span>
            </v-card-title>
            
            <v-card-text>
                <v-list lines="two">
                    <v-list-item
                        :title="`${patient.firstName} ${patient.lastName}`"
                        subtitle="Full Name"
                        prepend-icon="mdi-account"
                    ></v-list-item>
                    
                    <v-divider inset></v-divider>

                    <v-list-item
                        :title="patient.oib"
                        subtitle="OIB"
                        prepend-icon="mdi-card-account-details-outline"
                    ></v-list-item>

                    <v-divider inset></v-divider>

                    <v-list-item
                        :title="new Date(patient.birthDate).toLocaleDateString('en-GB')"
                        subtitle="Date of Birth"
                        prepend-icon="mdi-calendar"
                    ></v-list-item>

                    <v-divider inset></v-divider>

                    <v-list-item
                        :title="patient.gender"
                        subtitle="Gender"
                        prepend-icon="mdi-gender-male-female"
                    ></v-list-item>
                </v-list>
            </v-card-text>
        </v-card>

        <v-alert v-else-if="!isLoading" type="error" class="mt-4">
            Patient not found.
        </v-alert>

        <div v-if="isLoading" class="text-center pa-10">
            <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </div>
    </v-container>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { usePatientStore } from '@/stores/patientStore';
import type { Patient } from '@/stores/patientStore';

const router = useRouter();
const patientStore = usePatientStore();

const patient = ref<Patient | null>(null);
const isLoading = ref(true);

onMounted(() => {
  if (patientStore.selectedPatient) {
    patient.value = patientStore.selectedPatient;
  } else {
    router.push({ name: 'patient-list' });
  }
  isLoading.value = false;
});

function goBack() {
  patientStore.clearPatient();
  router.back();
}
</script>