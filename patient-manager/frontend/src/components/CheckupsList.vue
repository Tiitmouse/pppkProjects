<template>
    <v-card elevation="2">
        <v-card-title class="d-flex align-center justify-space-between">
            <span>Checkups</span>
            <v-btn @click="isCheckupDialogOpen = true" size="small" color="success" variant="elevated">
                <v-icon start>mdi-plus</v-icon>
                Add
            </v-btn>
        </v-card-title>
        <v-card-text class="pa-0">
            <div v-if="isLoadingCheckups" class="text-center pa-4">
                <v-progress-circular indeterminate color="primary"></v-progress-circular>
            </div>
            <v-data-table
                v-else
                :headers="checkupHeaders"
                :items="checkups"
                :group-by="[{ key: 'type', order: 'asc' }]"
                density="compact"
                item-value="uuid"
            >
                <template v-slot:item.checkupDate="{ item }">
                    {{ new Date(item.checkupDate).toLocaleDateString('hr-HR') }}
                </template>

                <template v-slot:item.checkupTime="{ item }">
                    {{ new Date(item.checkupDate).toLocaleTimeString('hr-HR', { hour: '2-digit', minute: '2-digit' }) }}
                </template>

                <template v-slot:group-header="{ item, columns, toggleGroup, isGroupOpen }">
                    <tr>
                        <td :colspan="columns.length">
                            <VBtn
                                size="small"
                                variant="text"
                                :icon="isGroupOpen(item) ? '$expand' : '$next'"
                                @click="toggleGroup(item)"
                            ></VBtn>
                            <span class="font-weight-bold">{{ getFullCheckupTypeName(item.value) }}</span>
                        </td>
                    </tr>
                </template>

                <template v-slot:no-data>
                    <div class="text-center text-grey py-4">
                        No checkups recorded.
                    </div>
                </template>
            </v-data-table>
        </v-card-text>
    </v-card>

    <CheckupDialog
        v-if="patient"
        v-model="isCheckupDialogOpen"
        :patient="patient"
        @save="handleCreateCheckup"
    />
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import type { PropType } from 'vue';
import { getCheckupsForRecord, createCheckup } from '@/services/patientService';
import type { Patient } from '@/stores/patientStore';
import type { CheckupDto, CreateCheckupDto } from '@/dtos/checkupDto';
import { CheckupType } from '@/enums/checkupType';
import CheckupDialog from '@/components/checkupDialog.vue';

const props = defineProps({
    patient: {
        type: Object as PropType<Patient | null>,
        required: true,
    }
});

const emit = defineEmits<{
  (e: 'show-snackbar', text: string, color: 'success' | 'error'): void
}>();

const checkups = ref<CheckupDto[]>([]);
const isLoadingCheckups = ref(true);
const isCheckupDialogOpen = ref(false);

const checkupHeaders = [
    { title: 'Date', key: 'checkupDate', align: 'start' },
    { title: 'Time', key: 'checkupTime', align: 'start', sortable: false },
    { title: 'Associated Illness ID', key: 'IllnessID', align: 'end' },
] as const;

function getFullCheckupTypeName(typeValue: CheckupType): string {
    const typeKey = (Object.keys(CheckupType) as Array<keyof typeof CheckupType>).find(key => CheckupType[key] === typeValue);
    if (typeKey) {
        return typeKey.replace(/([A-Z])/g, ' $1').trim();
    }
    return typeValue;
}

async function loadCheckups() {
    if (props.patient?.medicalRecordUuid) {
        isLoadingCheckups.value = true;
        try {
            checkups.value = await getCheckupsForRecord(props.patient.medicalRecordUuid);
        } catch (error) {
            console.error("Failed to load checkups:", error);
            emit('show-snackbar', 'Failed to load checkups.', 'error');
        } finally {
            isLoadingCheckups.value = false;
        }
    }
}

async function handleCreateCheckup(checkupData: CreateCheckupDto) {
    try {
        await createCheckup(checkupData);
        emit('show-snackbar', 'Checkup added successfully.', 'success');
        isCheckupDialogOpen.value = false;
        await loadCheckups();
    } catch (error) {
        console.error("Failed to add checkup:", error);
        emit('show-snackbar', 'Failed to add checkup.', 'error');
    }
}

onMounted(() => {
    loadCheckups();
});
</script>