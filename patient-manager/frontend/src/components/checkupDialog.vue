<template>
    <v-dialog :model-value="modelValue" @update:model-value="closeDialog" persistent max-width="600px">
        <v-card>
            <v-card-title>
                <span class="text-h5">Add New Checkup</span>
            </v-card-title>
            <v-card-text>
                <v-form ref="form" v-model="isFormValid">
                    <v-container>
                        <v-row>
                            <v-col cols="12" sm="6">
                                <v-text-field v-model="checkupData.checkupDate" label="Checkup Date" type="date"
                                    :rules="[rules.required]" required></v-text-field>
                            </v-col>
                            <v-col cols="12" sm="6">
                                <v-text-field v-model="checkupData.checkupTime" label="Checkup Time" type="time"
                                    :rules="[rules.required]" required></v-text-field>
                            </v-col>
                            <v-col cols="12">
                                <v-select v-model="checkupData.type" :items="checkupTypes" item-title="text"
                                    item-value="value" label="Checkup Type" :rules="[rules.required]"
                                    required></v-select>
                            </v-col>
                            <v-col cols="12">
                                <v-text-field v-model.number="checkupData.illnessId"
                                    label="Associated Illness ID (Optional)" type="number" clearable></v-text-field>
                            </v-col>
                            <v-col cols="12">
                                <v-file-input v-model="checkupData.files" label="Upload Images (Optional)"
                                    prepend-icon="mdi-camera" multiple accept="image/*" clearable></v-file-input>
                            </v-col>
                        </v-row>
                    </v-container>
                </v-form>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="error" variant="elevated" @click="closeDialog">
                    Cancel
                </v-btn>
                <v-btn color="primary" variant="elevated" @click="saveCheckup" :disabled="!isFormValid">
                    Save Checkup
                </v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, watch, computed } from 'vue';
import type { PropType } from 'vue';
import type { Patient } from '@/stores/patientStore';
import { CheckupType } from '@/enums/checkupType';
import type { CreateCheckupDto } from '@/dtos/checkupDto';
import { createCheckup, uploadCheckupImages } from '@/services/patientService';

const props = defineProps({
    modelValue: {
        type: Boolean,
        required: true,
    },
    patient: {
        type: Object as PropType<Patient | null>,
        required: true,
    }
});

const emit = defineEmits(['update:modelValue', 'save']);

const form = ref<any>(null);
const isFormValid = ref(false);

const initialData = {
    checkupDate: new Date().toISOString().split('T')[0],
    checkupTime: new Date().toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false }),
    type: undefined,
    illnessId: undefined,
    files: [],
};

const checkupData = reactive({ ...initialData });

watch(() => props.modelValue, (isOpen) => {
    if (isOpen) {
        Object.assign(checkupData, initialData);
        form.value?.resetValidation();
    }
});

const checkupTypes = computed(() => {
    return Object.entries(CheckupType).map(([key, value]) => ({
        text: key.replace(/([A-Z])/g, ' $1').trim(),
        value: value,
    }));
});


const rules = {
    required: (value: any) => !!value || 'This field is required.',
};

function closeDialog() {
    emit('update:modelValue', false);
}

async function saveCheckup() {
    const { valid } = await form.value?.validate();
    if (!valid || !props.patient) return;

    const combinedDateTime = `${checkupData.checkupDate}T${checkupData.checkupTime}`;

    const payload: CreateCheckupDto = {
        checkupDate: new Date(combinedDateTime).toISOString(),
        type: checkupData.type!,
        medicalRecordUuid: props.patient.medicalRecordUuid,
    };

    if (checkupData.illnessId) {
        payload.illnessId = checkupData.illnessId;
    }

    try {
        const createdCheckup = await createCheckup(payload);
        if (checkupData.files && checkupData.files.length > 0) {
            await uploadCheckupImages(createdCheckup.uuid, checkupData.files);
        }
        emit('save'); 
    } catch (error) {
        console.error('Failed to create checkup:', error);
    }
}
</script>