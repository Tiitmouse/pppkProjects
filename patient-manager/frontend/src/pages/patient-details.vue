<template>
    <v-container>
        <div class="d-flex align-center mb-4">
            <v-btn @click="goBack" color="grey">
                <v-icon start>mdi-arrow-left</v-icon>
                Back to list
            </v-btn>

            <v-spacer></v-spacer>

            <v-btn @click="editPatient" color="primary" class="ml-4">
                <v-icon start>mdi-pencil</v-icon>
                Edit
            </v-btn>
            <v-btn @click="confirmAndDelete" color="error" class="ml-4" :loading="isDeleting">
                <v-icon start>mdi-trash-can</v-icon>
                Delete
            </v-btn>
            <v-btn @click="openDoctorDialog" color="secondary" class="ml-4">
                <v-icon start>mdi-doctor</v-icon>
                Doctor
            </v-btn>
        </div>

        <v-card v-if="patient" :loading="isLoading">
             <v-card-title class="pa-4">
                <span class="text-h5">Patient Details</span>
            </v-card-title>
            <v-card-text>
                <v-list lines="two">
                    <v-list-item :title="`${patient.firstName} ${patient.lastName}`" subtitle="Full Name"
                        prepend-icon="mdi-account"></v-list-item>
                    <v-divider inset></v-divider>
                    <v-list-item :title="patient.oib" subtitle="OIB"
                        prepend-icon="mdi-card-account-details-outline"></v-list-item>
                    <v-divider inset></v-divider>
                    <v-list-item :title="new Date(patient.birthDate).toLocaleDateString('hr-HR')"
                        subtitle="Date of Birth" prepend-icon="mdi-calendar"></v-list-item>
                    <v-divider inset></v-divider>
                    <v-list-item :title="patient.gender" subtitle="Gender"
                        prepend-icon="mdi-gender-male-female"></v-list-item>
                </v-list>
            </v-card-text>
        </v-card>

        <v-alert v-else-if="!isLoading" type="error" class="mt-4">
            Patient not found.
        </v-alert>
        <div v-if="isLoading" class="text-center pa-10">
            <v-progress-circular indeterminate color="primary"></v-progress-circular>
        </div>

        <ConfirmDialogue ref="confirmDialog" />
        
        <v-dialog v-model="isDoctorDialogVisible" max-width="400px">
            </v-dialog>
        
        <v-snackbar v-model="snackbar.visible" :color="snackbar.color" :timeout="3000">
            {{ snackbar.text }}
        </v-snackbar>

    </v-container>
</template>

<script lang="ts" setup>
import { ref, onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { usePatientStore } from '@/stores/patientStore';
import type { Patient } from '@/stores/patientStore';
import { deletePatient } from '@/services/patientService';
import ConfirmDialogue from '@/components/confirmDialog.vue';

const router = useRouter();
const patientStore = usePatientStore();

const patient = ref<Patient | null>(null);
const isLoading = ref(true);
const isDeleting = ref(false);

const confirmDialog = ref();
const isDoctorDialogVisible = ref(false);
const mockDoctor = ref({ firstName: 'Ana', lastName: 'Anić' });
const snackbar = reactive({
    visible: false,
    text: '',
    color: 'success',
});


onMounted(() => {
    if (patientStore.selectedPatient) {
        patient.value = patientStore.selectedPatient;
    } else {
        router.push({ name: 'patient-list' });
    }
    isLoading.value = false;
});

function showSnackbar(text: string, color: 'success' | 'error') {
    snackbar.text = text;
    snackbar.color = color;
    snackbar.visible = true;
}

function goBack() {
    patientStore.clearPatient();
    router.back();
}

function editPatient() {
    alert('Edit action clicked. (Not implemented yet)');
}

async function confirmAndDelete() {
    if (!patient.value) return;

    const isConfirmed = await confirmDialog.value.Open({
        Title: 'Delete Patient',
        Message: `Are you sure you want to delete ${patient.value.firstName} ${patient.value.lastName}? This action cannot be undone.`,
    });

    if (isConfirmed) {
        isDeleting.value = true;
        try {
            await deletePatient(patient.value.id);
            showSnackbar('Patient deleted successfully.', 'success');
            router.push({ name: 'patient-list' });
        } catch (error) {
            console.error("Failed to delete patient:", error);
            showSnackbar('Failed to delete patient. Please try again.', 'error');
        } finally {
            isDeleting.value = false;
        }
    }
}

function openDoctorDialog() {
    isDoctorDialogVisible.value = true;
}
</script>