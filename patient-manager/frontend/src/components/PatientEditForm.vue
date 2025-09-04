<template>
    <v-form @submit.prevent="submit">
        <v-container>
            <v-row>
                <v-col cols="12" md="6">
                    <v-text-field
                        v-model="editablePatient.firstName"
                        label="First Name"
                        variant="outlined"
                        density="compact"
                        required
                    ></v-text-field>
                </v-col>
                <v-col cols="12" md="6">
                    <v-text-field
                        v-model="editablePatient.lastName"
                        label="Last Name"
                        variant="outlined"
                        density="compact"
                        required
                    ></v-text-field>
                </v-col>
                <v-col cols="12">
                    <v-text-field
                        v-model="editablePatient.oib"
                        label="OIB"
                        variant="outlined"
                        density="compact"
                        required
                    ></v-text-field>
                </v-col>
                <v-col cols="12">
                     <v-text-field
                        v-model="editablePatient.birthDate"
                        label="Date of Birth"
                        type="date"
                        variant="outlined"
                        density="compact"
                        required
                    ></v-text-field>
                </v-col>
            </v-row>
             <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="error" variant="elevated" @click="$emit('cancel')">Cancel</v-btn>
                <v-btn color="primary" variant="elevated" type="submit" :loading="isSaving">Save</v-btn>
            </v-card-actions>
        </v-container>
    </v-form>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import type { Patient } from '@/stores/patientStore';

const props = defineProps<{
  patient: Patient | null,
  isSaving: boolean
}>();

const emit = defineEmits<{
  (e: 'save', patientData: Patient): void,
  (e: 'cancel'): void
}>();

const editablePatient = ref<Patient>({} as Patient);

watch(() => props.patient, (newVal) => {
  if (newVal) {
    const date = new Date(newVal.birthDate);
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    
    editablePatient.value = { 
        ...newVal,
        birthDate: `${year}-${month}-${day}`
    };
  }
}, { immediate: true });

function submit() {
    const patientDataToSave = {
        ...editablePatient.value,
        birthDate: new Date(editablePatient.value.birthDate).toISOString()
    };
    emit('save', patientDataToSave);
}
</script>