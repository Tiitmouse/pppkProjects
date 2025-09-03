<template>
    <v-container class="mt-4">
        <v-text-field
            v-model="search"
            label="Search patients"
            prepend-inner-icon="mdi-magnify"
            variant="outlined"
            hide-details
            single-line
            class="mb-4"
        ></v-text-field>

        <v-data-table
            :headers="headers"
            :items="filteredPatients"
            item-value="id"
            fixed-header
            height="60vh"
        >
            <template v-slot:no-data>
                No patients found.
            </template>
            
            <template v-slot:item.birthDate="{ item }">
                {{ new Date(item.birthDate).toLocaleDateString('en-GB') }}
            </template>

            <template v-slot:item.actions="{ item }">
                <div class="d-flex justify-end">
                    <v-icon 
                        color="medium-emphasis" 
                        icon="mdi-eye" 
                        size="large" 
                        @click="navigateToPatientDetails(item)"
                    ></v-icon>
                </div>
            </template>
        </v-data-table>
    </v-container>
</template>

<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import * as patientService from '@/services/patientService';
import { usePatientStore } from '@/stores/patientStore';
import type { Patient } from '@/stores/patientStore';

const router = useRouter();
const patientStore = usePatientStore();

const headers = [
    { title: 'First Name', key: 'firstName', align: 'start' },
    { title: 'Last Name', key: 'lastName' },
    { title: 'OIB', key: 'oib' },
    { title: 'Date of Birth', key: 'birthDate' },
    { title: 'Gender', key: 'gender' },
    { title: 'Details', key: 'actions', sortable: false, align: 'end' },
] as const;

const patients = ref<Patient[]>([]);
const search = ref('');

const filteredPatients = computed(() => {
    if (!search.value) {
        return patients.value;
    }
    const searchTerm = search.value.toLowerCase();
    return patients.value.filter(patient => 
        patient.firstName.toLowerCase().includes(searchTerm) ||
        patient.lastName.toLowerCase().includes(searchTerm) ||
        patient.oib.toLowerCase().includes(searchTerm)
    );
});

async function loadPatients() {
    try {
        const data = await patientService.getAllPatients();
        patients.value = data || [];
    } catch (error) {
        console.error("Error fetching patients:", error);
    }
}

function navigateToPatientDetails(patient: Patient) {
    patientStore.viewPatient(patient); 
    router.push({ name: 'patient-details' });
}

onMounted(() => {
    loadPatients();
});
</script>