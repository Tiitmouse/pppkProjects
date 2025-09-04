<template>
    <v-container>
        <v-card elevation="7" variant="tonal" class="round-xl pa-4 mb-5">
            <div class="d-flex align-center mb-4">
                <v-btn @click="goBack" color="grey">
                    <v-icon start>mdi-arrow-left</v-icon>
                    Back to list
                </v-btn>

                <v-spacer></v-spacer>

                <v-btn @click="openEditOptions" color="primary" class="ml-4">
                    <v-icon start>mdi-pencil</v-icon>
                    Edit
                </v-btn>
                <v-btn @click="confirmAndDelete" color="error" class="ml-4" :loading="isDeleting">
                    <v-icon start>mdi-trash-can</v-icon>
                    Delete
                </v-btn>
                <v-btn @click="isDoctorDialogVisible = true" color="secondary" class="ml-4">
                    <v-icon start>mdi-doctor</v-icon>
                    Doctor
                </v-btn>
            </div>

            <v-card v-if="patient" elevation="7" variant="tonal" class="round-xl pa-4">
            <div v-if="patient" :loading="isLoading">
                <v-card-title class="pa-0 pt-2 pb-2">
                    <span class="text-h5">PATIENT CARTON</span>
                </v-card-title>
                <v-card-text class="pa-0">
                    <v-list lines="two" class="bg-transparent">
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
            </div>
            </v-card>

            <v-alert v-else-if="!isLoading" type="error" class="mt-4">
                Patient not found.
            </v-alert>

            <div v-if="isLoading" class="text-center pa-10">
                <v-progress-circular indeterminate color="primary"></v-progress-circular>
            </div>
        </v-card>
    <v-card v-if="patient" elevation="7" variant="tonal" class="round-xl pa-4">
             <v-row>
                <v-col cols="12" md="6">
                    <IllnessesList
                        :patient="patient"
                        :is-editing="illnessesInEditMode"
                        @show-snackbar="showSnackbar"
                    />
                </v-col>
                <v-col cols="12" md="6">
                    <CheckupsList
                        :patient="patient"
                        :is-editing="checkupsInEditMode"
                        @show-snackbar="showSnackbar"
                    />
                </v-col>
            </v-row>
        </v-card>

        <OptionsDialogue ref="editOptionsDialog" />
        <ConfirmDialogue ref="confirmDialog" />

        <v-dialog v-model="isEditDialogOpen" persistent max-width="600px">
            <v-card>
                 <v-card-title>
                    <span class="text-h5">Edit Patient Data</span>
                </v-card-title>
                <v-card-text>
                    <PatientEditForm
                        :patient="patient"
                        :is-saving="isSaving"
                        @save="handleUpdatePatient"
                        @cancel="isEditDialogOpen = false"
                    />
                </v-card-text>
            </v-card>
        </v-dialog>

        <v-dialog v-model="isDoctorDialogVisible" max-width="400px">
            <v-card>
                <v-card-title>
                    <span class="text-h5">Doctor Assignment</span>
                </v-card-title>
                <v-card-text>
                    This patient is assigned to:
                    <div class="font-weight-bold mt-2">
                        Dr. {{ mockDoctor.firstName }} {{ mockDoctor.lastName }}
                    </div>
                </v-card-text>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="primary" variant="elevated" @click="isDoctorDialogVisible = false">
                        OK
                    </v-btn>
                </v-card-actions>
            </v-card>
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
import { deletePatient, updatePatient } from '@/services/patientService';
import OptionsDialogue from '@/components/optionsDialog.vue';
import ConfirmDialogue from '@/components/confirmDialog.vue';
import PatientEditForm from '@/components/PatientEditForm.vue';
import CheckupsList from '@/components/CheckupsList.vue';
import IllnessesList from '@/components/IllnessesList.vue';

const router = useRouter();
const patientStore = usePatientStore();

const patient = ref<Patient | null>(null);
const isLoading = ref(true);
const isDeleting = ref(false);
const isSaving = ref(false);
const checkupsInEditMode = ref(false);
const illnessesInEditMode = ref(false);

const editOptionsDialog = ref();
const confirmDialog = ref();
const isDoctorDialogVisible = ref(false);
const isEditDialogOpen = ref(false);

const mockDoctor = ref({ firstName: 'Ana', lastName: 'Anić' });
const snackbar = reactive({
    visible: false,
    text: '',
    color: 'success' as 'success' | 'error' | 'info',
});

onMounted(() => {
    if (patientStore.selectedPatient) {
        patient.value = patientStore.selectedPatient;
    } else {
        router.push({ name: 'patient-list' });
    }
    isLoading.value = false;
});

function showSnackbar(text: string, color: 'success' | 'error' | 'info') {
    snackbar.text = text;
    snackbar.color = color;
    snackbar.visible = true;
}

function goBack() {
    patientStore.clearPatient();
    router.back();
}

async function openEditOptions() {
    const selectedChoice = await editOptionsDialog.value.Open({
        Title: 'What would you like to edit?',
        Options: ['Patient Data', 'Checkups', 'Illnesses', 'Prescriptions']
    });

    if (selectedChoice === 'Patient Data') {
        isEditDialogOpen.value = true;
    } else if (selectedChoice === 'Checkups') {
        checkupsInEditMode.value = !checkupsInEditMode.value;
        const status = checkupsInEditMode.value ? 'enabled' : 'disabled';
        showSnackbar(`Checkups editing ${status}.`, 'info');
    } else if (selectedChoice === 'Illnesses') {
        illnessesInEditMode.value = !illnessesInEditMode.value;
        const status = illnessesInEditMode.value ? 'enabled' : 'disabled';
        showSnackbar(`Illnesses editing ${status}.`, 'info');
    } else if (selectedChoice === 'Prescriptions') {
        await confirmDialog.value.Open({
            Title: 'Manage Prescriptions',
            Message: `To manage prescriptions, please use the "View" button next to an illness in the illnesses list.`,
            Options: { noCancel: true }
        });
    }
}

async function handleUpdatePatient(updatedPatientData: Patient) {
    isSaving.value = true;
    try {
        const updatedPatient = await updatePatient(updatedPatientData.id, updatedPatientData);
        patient.value = updatedPatient;
        patientStore.selectedPatient = updatedPatient;
        isEditDialogOpen.value = false;
        showSnackbar('Patient updated successfully.', 'success');
    } catch (error) {
        console.error("Failed to update patient:", error);
        showSnackbar('Failed to update patient.', 'error');
    } finally {
        isSaving.value = false;
    }
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
</script>