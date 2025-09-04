<template>
    <v-card elevation="2">
        <v-card-title class="d-flex align-center justify-space-between">
            <span>Illnesses</span>
            <v-btn @click="openCreateDialog" size="small" color="success" variant="elevated">
                <v-icon start>mdi-plus</v-icon>
                Add
            </v-btn>
        </v-card-title>
        <v-card-text class="pa-0">
            <div v-if="isLoading" class="text-center pa-4">
                <v-progress-circular indeterminate color="primary"></v-progress-circular>
            </div>
            <v-data-table
                v-else
                :headers="headers"
                :items="illnesses"
                density="compact"
                item-value="id"
            >
                <template v-slot:item.startDate="{ item }">
                    {{ new Date(item.startDate).toLocaleDateString('hr-HR') }}
                </template>
                <template v-slot:item.endDate="{ item }">
                    {{ item.endDate ? new Date(item.endDate).toLocaleDateString('hr-HR') : 'Ongoing' }}
                </template>
                 <template v-slot:item.actions="{ item }">
                    <div class="d-flex justify-end">
                        <v-btn icon @click="$emit('view-prescriptions', item)" size="small" variant="tonal" color="info" class="me-2">
                            <v-icon>mdi-pill</v-icon>
                        </v-btn>
                        <template v-if="isEditing">
                            <v-icon size="small" class="me-2" @click="openEditDialog(item)">mdi-pencil</v-icon>
                            <v-icon size="small" @click="confirmAndDelete(item)">mdi-delete</v-icon>
                        </template>
                    </div>
                </template>
                <template v-slot:no-data>
                    <div class="text-center text-grey py-4">No illnesses recorded.</div>
                </template>
            </v-data-table>
        </v-card-text>
    </v-card>

    <v-dialog v-model="isDialogOpen" persistent max-width="600px">
        <v-card>
            <v-card-title>
                <span class="text-h5">{{ isEditingDialog ? 'Edit' : 'Add' }} Illness</span>
            </v-card-title>
            <v-card-text>
                <v-form ref="form" v-model="isFormValid">
                    <v-text-field v-model="formData.name" label="Illness Name" :rules="[rules.required]"></v-text-field>
                    <v-text-field v-model="formData.startDate" label="Start Date" type="date" :rules="[rules.required]"></v-text-field>
                    <v-text-field v-model="formData.endDate" label="End Date (Optional)" type="date" clearable></v-text-field>
                </v-form>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="error" variant="text" @click="closeDialog">Cancel</v-btn>
                <v-btn color="primary" variant="elevated" @click="save" :disabled="!isFormValid">Save</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>

    <ConfirmDialogue ref="confirmDialog" />
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch } from 'vue';
import type { PropType } from 'vue';
import type { Patient } from '@/stores/patientStore';
import ConfirmDialogue from '@/components/confirmDialog.vue';
import type { IllnessListDto, UpdateIllnessDto, CreateIllnessDto } from '@/dtos/illnessDto';
import { getIllnessesForRecord, updateIllness, createIllness, deleteIllness } from '@/services/patientService';

const props = defineProps({
    patient: { type: Object as PropType<Patient | null>, required: true },
    isEditing: { type: Boolean, default: false },
});

const emit = defineEmits(['show-snackbar', 'view-prescriptions']);

const illnesses = ref<IllnessListDto[]>([]);
const isLoading = ref(true);
const isDialogOpen = ref(false);
const isEditingDialog = ref(false);
const isFormValid = ref(false);
const form = ref<any>(null);
const selectedIllness = ref<IllnessListDto | null>(null);
const confirmDialog = ref();

const formData = reactive({ name: '', startDate: '', endDate: '' });

const rules = { required: (v: any) => !!v || 'This field is required.' };

const headers = computed(() => {
    const baseHeaders: readonly {
        title: string;
        key: string;
        align?: 'start' | 'center' | 'end';
        width?: string;
        sortable?: boolean;
    }[] = [
        { title: 'ID', key: 'id', align: 'start', width: '15%' },
        { title: 'Name', key: 'name', align: 'start' },
        { title: 'Start Date', key: 'startDate' },
        { title: 'End Date', key: 'endDate' },
    ];

    const actionsHeader: {
        title: string;
        key: string;
        sortable: boolean;
        align: 'start' | 'center' | 'end';
    } = { title: 'Actions', key: 'actions', sortable: false, align: 'end' };

    if (props.isEditing) {
        return [...baseHeaders, actionsHeader];
    }
    return baseHeaders.concat(actionsHeader);
});


async function loadIllnesses() {
    if (!props.patient) return;
    isLoading.value = true;
    try {
        illnesses.value = await getIllnessesForRecord(props.patient.medicalRecordUuid);
    } catch (error) {
        emit('show-snackbar', 'Failed to load illnesses.', 'error');
    } finally {
        isLoading.value = false;
    }
}

function openCreateDialog() {
    isEditingDialog.value = false;
    Object.assign(formData, { name: '', startDate: new Date().toISOString().split('T')[0], endDate: '' });
    isDialogOpen.value = true;
}

function openEditDialog(illness: IllnessListDto) {
    isEditingDialog.value = true;
    selectedIllness.value = illness;
    Object.assign(formData, {
        name: illness.name,
        startDate: new Date(illness.startDate).toISOString().split('T')[0],
        endDate: illness.endDate ? new Date(illness.endDate).toISOString().split('T')[0] : ''
    });
    isDialogOpen.value = true;
}

function closeDialog() {
    isDialogOpen.value = false;
    form.value?.resetValidation();
}

async function save() {
    await form.value?.validate();
    if (!isFormValid.value || !props.patient) return;

    const data = {
        name: formData.name,
        startDate: new Date(formData.startDate).toISOString(),
        endDate: formData.endDate ? new Date(formData.endDate).toISOString() : undefined,
    };

    try {
        if (isEditingDialog.value && selectedIllness.value) {
            await updateIllness(selectedIllness.value.uuid, data as UpdateIllnessDto);
            emit('show-snackbar', 'Illness updated successfully.', 'success');
        } else {
            const createData: CreateIllnessDto = { ...data, medicalRecordUuid: props.patient.medicalRecordUuid };
            await createIllness(createData);
            emit('show-snackbar', 'Illness added successfully.', 'success');
        }
        closeDialog();
        await loadIllnesses();
    } catch (error) {
        emit('show-snackbar', `Failed to save illness.`, 'error');
    }
}

async function confirmAndDelete(illness: IllnessListDto) {
    const isConfirmed = await confirmDialog.value.Open({
        Title: 'Delete Illness',
        Message: `Are you sure you want to delete the illness: "${illness.name}"?`,
    });
    if (isConfirmed) {
        try {
            await deleteIllness(illness.uuid);
            emit('show-snackbar', 'Illness deleted successfully.', 'success');
            await loadIllnesses();
        } catch (error) {
            emit('show-snackbar', 'Failed to delete illness.', 'error');
        }
    }
}

watch(() => props.patient, (newPatient) => {
    if (newPatient) {
        loadIllnesses();
    }
}, { immediate: true });
</script>