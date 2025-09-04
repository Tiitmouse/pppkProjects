<template>
    <v-card>
        <v-card-title class="d-flex align-center justify-space-between">
            <span>Prescriptions for "{{ illness.name }}"</span>
            <v-btn @click="isDialogOpen = true" size="small" color="success" variant="elevated">
                <v-icon start>mdi-plus</v-icon>
                Add
            </v-btn>
        </v-card-title>
        <v-card-text>
            <div v-if="isLoading" class="text-center pa-4">
                <v-progress-circular indeterminate color="primary"></v-progress-circular>
            </div>
            <v-data-table v-else :headers="headers" :items="prescriptions" density="compact">
                <template v-slot:item.issuedAt="{ item }">
                    {{ new Date(item.issuedAt).toLocaleDateString('hr-HR') }}
                </template>
                <template v-slot:item.medications="{ item }">
                    <v-chip v-for="med in item.medications" :key="med.uuid" size="small" class="mr-1 mb-1">
                        {{ med.name }}
                    </v-chip>
                </template>
                <template v-slot:item.actions="{ item }">
                     <div class="d-flex justify-end">
                        <v-icon size="small" @click="confirmAndDelete(item)">mdi-delete</v-icon>
                    </div>
                </template>
                 <template v-slot:no-data>
                    <div class="text-center text-grey py-4">No prescriptions found for this illness.</div>
                </template>
            </v-data-table>
        </v-card-text>
    </v-card>

    <v-dialog v-model="isDialogOpen" persistent max-width="600px">
        <v-card>
            <v-card-title>Add Prescription</v-card-title>
            <v-card-text>
                <v-form ref="form" v-model="isFormValid">
                    <v-text-field v-model="formData.issuedAt" label="Date Issued" type="date" :rules="[rules.required]"></v-text-field>
                    <v-autocomplete
                        v-model="formData.medicationUuids"
                        :items="allMedications"
                        item-title="name"
                        item-value="uuid"
                        label="Medications"
                        multiple
                        chips
                        closable-chips
                    ></v-autocomplete>
                </v-form>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="error" variant="text" @click="isDialogOpen = false">Cancel</v-btn>
                <v-btn color="primary" variant="elevated" @click="save" :disabled="!isFormValid">Save</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>

    <ConfirmDialogue ref="confirmDialog" />
</template>

<script lang="ts" setup>
import { ref, onMounted, reactive } from 'vue';
import type { PropType } from 'vue';
import ConfirmDialogue from '@/components/confirmDialog.vue';
import type { IllnessListDto } from '@/dtos/illnessDto';
import type { PrescriptionListDto, MedicationListDto, CreatePrescriptionDto } from '@/dtos/prescriptionDto';
import { getPrescriptionsForIllness, getAllMedications, createPrescription, deletePrescription } from '@/services/patientService';

const props = defineProps({
    illness: { type: Object as PropType<IllnessListDto>, required: true },
});

const emit = defineEmits(['show-snackbar']);

const prescriptions = ref<PrescriptionListDto[]>([]);
const allMedications = ref<MedicationListDto[]>([]);
const isLoading = ref(true);
const isDialogOpen = ref(false);
const isFormValid = ref(false);
const form = ref<any>(null);
const confirmDialog = ref();

const formData = reactive({
    issuedAt: new Date().toISOString().split('T')[0],
    medicationUuids: [],
});

const rules = { required: (v: any) => !!v || 'This field is required.' };

const headers = [
    { title: 'Date Issued', key: 'issuedAt', align: 'start' },
    { title: 'Medications', key: 'medications' },
    { title: 'Actions', key: 'actions', sortable: false, align: 'end' },
] as const;

async function loadData() {
    isLoading.value = true;
    try {
        const [prescriptionsData, medicationsData] = await Promise.all([
            getPrescriptionsForIllness(props.illness.id),
            getAllMedications()
        ]);
        prescriptions.value = prescriptionsData;
        allMedications.value = medicationsData;
    } catch (error) {
        emit('show-snackbar', 'Failed to load prescription data.', 'error');
    } finally {
        isLoading.value = false;
    }
}

async function save() {
    await form.value?.validate();
    if (!isFormValid.value) return;

    const payload: CreatePrescriptionDto = {
        illnessId: props.illness.id,
        issuedAt: new Date(formData.issuedAt).toISOString(),
        medicationUuids: formData.medicationUuids,
    };

    try {
        await createPrescription(payload);
        emit('show-snackbar', 'Prescription added successfully.', 'success');
        isDialogOpen.value = false;
        await loadData();
    } catch (error) {
        emit('show-snackbar', 'Failed to add prescription.', 'error');
    }
}

async function confirmAndDelete(prescription: PrescriptionListDto) {
    const isConfirmed = await confirmDialog.value.Open({
        Title: 'Delete Prescription',
        Message: `Are you sure you want to delete this prescription?`,
    });
    if (isConfirmed) {
        try {
            await deletePrescription(prescription.uuid);
            emit('show-snackbar', 'Prescription deleted successfully.', 'success');
            await loadData();
        } catch (error) {
            emit('show-snackbar', 'Failed to delete prescription.', 'error');
        }
    }
}

onMounted(loadData);
</script>