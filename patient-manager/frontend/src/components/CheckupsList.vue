<template>
    <v-card elevation="2">
        <v-card-title class="d-flex align-center justify-space-between">
            <span>Checkups</span>
            <v-btn @click="isCreateDialogOpen = true" size="small" color="success" variant="elevated">
                <v-icon start>mdi-plus</v-icon>
                Add
            </v-btn>
        </v-card-title>
        <v-card-text class="pa-0">
            <div v-if="isLoadingCheckups" class="text-center pa-4">
                <v-progress-circular indeterminate color="primary"></v-progress-circular>
            </div>
            <v-data-table v-else :headers="checkupHeaders" :items="checkups" :group-by="[{ key: 'type', order: 'asc' }]"
                density="compact" item-value="uuid">
                <template v-slot:item.checkupDate="{ item }">
                    {{ new Date(item.checkupDate).toLocaleDateString('hr-HR') }}
                </template>

                <template v-slot:item.checkupTime="{ item }">
                    {{ new Date(item.checkupDate).toLocaleTimeString('hr-HR', { hour: '2-digit', minute: '2-digit' }) }}
                </template>

                <template v-slot:item.actions="{ item }">
                    <div class="d-flex justify-end">
                        <v-icon size="small" class="me-2" @click="openEditDialog(item)">mdi-pencil</v-icon>
                        <v-icon size="small" @click="confirmAndDelete(item)">mdi-delete</v-icon>
                    </div>
                </template>

                <template v-slot:item.images="{ item }">
                    <div v-if="item.images && item.images.length > 0">
                        <v-btn icon @click="openGallery(item.images)" size="small" variant="tonal" color="info"
                            class="me-2">
                            <v-icon color="info">mdi-image-multiple-outline</v-icon>
                        </v-btn>
                    </div>
                    <div v-else>
                        <v-icon color="info">mdi-image-off-outline</v-icon>
                    </div>
                </template>

                <template v-slot:group-header="{ item, columns, toggleGroup, isGroupOpen }">
                    <tr>
                        <td :colspan="columns.length">
                            <VBtn size="small" variant="text" :icon="isGroupOpen(item) ? '$expand' : '$next'"
                                @click="toggleGroup(item)"></VBtn>
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

    <CheckupDialog v-if="patient" v-model="isCreateDialogOpen" :patient="patient" @save="handleCreateCheckup" />

    <v-dialog v-model="isEditDialogOpen" persistent max-width="600px">
        <v-card v-if="selectedCheckup">
            <v-card-title>
                <span class="text-h5">Edit Checkup</span>
            </v-card-title>
            <v-card-text>
                <v-form ref="editForm" v-model="isEditFormValid">
                    <v-container>
                        <v-row>
                            <v-col cols="12" sm="6">
                                <v-text-field v-model="editFormData.checkupDate" label="Checkup Date" type="date"
                                    :rules="[rules.required]" required></v-text-field>
                            </v-col>
                            <v-col cols="12" sm="6">
                                <v-text-field v-model="editFormData.checkupTime" label="Checkup Time" type="time"
                                    :rules="[rules.required]" required></v-text-field>
                            </v-col>
                            <v-col cols="12">
                                <v-select v-model="editFormData.type" :items="checkupTypes" item-title="text"
                                    item-value="value" label="Checkup Type" :rules="[rules.required]"
                                    required></v-select>
                            </v-col>
                            <v-col cols="12">
                                <v-text-field v-model.number="editFormData.illnessId"
                                    label="Associated Illness ID (Optional)" type="number" clearable></v-text-field>
                            </v-col>
                            <v-col cols="12">
                                <v-file-input v-model="newFiles" label="Add More Images" prepend-icon="mdi-camera"
                                    multiple accept="image/*" clearable></v-file-input>
                            </v-col>
                            <v-col v-if="existingImages.length > 0" cols="12">
                                <p class="mb-2 text-subtitle-1">Existing Images</p>
                                <div class="d-flex flex-wrap ga-2">
                                    <div v-for="image in existingImages" :key="image.uuid" class="d-flex flex-column align-center">
                                        <v-img :src="'http://localhost:8098/api/checkup/image/' + image.path" height="100" width="100" cover aspect-ratio="1" class="rounded">
                                            <template v-slot:placeholder>
                                                <div class="d-flex align-center justify-center fill-height">
                                                    <v-progress-circular indeterminate color="grey-lighten-4"></v-progress-circular>
                                                </div>
                                            </template>
                                        </v-img>
                                        <v-btn size="small" color="error" class="mt-1" variant="tonal" @click="deleteImage(image)">
                                            <v-icon>mdi-delete</v-icon>
                                        </v-btn>
                                    </div>
                                </div>
                            </v-col>
                        </v-row>
                    </v-container>
                </v-form>
            </v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="error" variant="text" @click="isEditDialogOpen = false">Cancel</v-btn>
                <v-btn color="primary" variant="elevated" @click="handleUpdateCheckup"
                    :disabled="!isEditFormValid">Save</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>

    <ConfirmDialogue ref="confirmDialog" />

    <v-dialog v-model="isGalleryOpen" max-width="800px">
        <v-card>
            <v-toolbar color="primary" density="compact">
                <v-toolbar-title>Image Gallery</v-toolbar-title>
                <v-spacer></v-spacer>
                <v-btn icon @click="isGalleryOpen = false">
                    <v-icon>mdi-close</v-icon>
                </v-btn>
            </v-toolbar>
            <v-carousel v-if="selectedCheckupImages.length > 0" show-arrows="hover" hide-delimiters>
                <v-carousel-item 
                    v-for="image in selectedCheckupImages" 
                    :key="image.uuid"
                    :src="'http://localhost:8098/api/checkup/image/' + image.path">
                </v-carousel-item>
            </v-carousel>
             <v-card-text v-else class="text-center pa-10">
                <p>No images available for this checkup.</p>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, watch } from 'vue';
import type { PropType } from 'vue';
import { getCheckupsForRecord, createCheckup, updateCheckup, deleteCheckup, uploadCheckupImages, deleteCheckupImage } from '@/services/patientService';
import type { Patient } from '@/stores/patientStore';
import { type CheckupDto, type CreateCheckupDto, type UpdateCheckupDto, type ImageDto } from '@/dtos/checkupDto';
import { CheckupType } from '@/enums/checkupType';
import CheckupDialog from '@/components/checkupDialog.vue';
import ConfirmDialogue from '@/components/confirmDialog.vue';

const props = defineProps({
    patient: {
        type: Object as PropType<Patient | null>,
        required: true,
    },
    isEditing: {
        type: Boolean,
        default: false,
    },
});

const emit = defineEmits<{
    (e: 'show-snackbar', text: string, color: 'success' | 'error' | 'info'): void
}>();

const checkups = ref<CheckupDto[]>([]);
const isLoadingCheckups = ref(true);
const isCreateDialogOpen = ref(false);
const isEditDialogOpen = ref(false);
const isEditFormValid = ref(false);
const editForm = ref<any>(null);
const selectedCheckup = ref<CheckupDto | null>(null);
const confirmDialog = ref();
const isGalleryOpen = ref(false);
const selectedCheckupImages = ref<ImageDto[]>([]);
const newFiles = ref<File[]>([]);
const existingImages = ref<ImageDto[]>([]);

const editFormData = reactive({
    checkupDate: '',
    checkupTime: '',
    type: '' as CheckupType | undefined,
    illnessId: undefined as number | undefined,
});

const baseHeaders = [
    { title: 'Date', key: 'checkupDate', align: 'start' },
    { title: 'Time', key: 'checkupTime', align: 'start', sortable: false },
    { title: 'Associated Illness ID', key: 'illnessId', align: 'end' },
    { title: 'Images', key: 'images', sortable: false, align: 'center' },
] as const;

const checkupHeaders = computed(() => {
    if (props.isEditing) {
        return [
            ...baseHeaders,
            { title: 'Actions', key: 'actions', sortable: false, align: 'end' },
        ] as const;
    }
    return baseHeaders;
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
            emit('show-snackbar', 'Failed to load checkups.', 'error');
        } finally {
            isLoadingCheckups.value = false;
        }
    }
}

async function handleCreateCheckup() {
    try {
        await loadCheckups();
        emit('show-snackbar', 'Checkup added successfully.', 'success');
        isCreateDialogOpen.value = false;
    } catch (error) {
        emit('show-snackbar', 'Failed to add checkup.', 'error');
    }
}

function openEditDialog(checkup: CheckupDto) {
    selectedCheckup.value = checkup;
    const date = new Date(checkup.checkupDate);
    editFormData.checkupDate = date.toISOString().split('T')[0];
    editFormData.checkupTime = date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: false });
    editFormData.type = checkup.type;
    editFormData.illnessId = checkup.illnessId;
    existingImages.value = checkup.images || [];
    newFiles.value = [];
    isEditDialogOpen.value = true;
}

async function handleUpdateCheckup() {
    const { valid } = await editForm.value?.validate();
    if (!valid || !selectedCheckup.value) return;

    const combinedDateTime = `${editFormData.checkupDate}T${editFormData.checkupTime}`;
    const payload: UpdateCheckupDto = {
        checkupDate: new Date(combinedDateTime).toISOString(),
        type: editFormData.type!,
        illnessId: editFormData.illnessId,
    };

    try {
        await updateCheckup(selectedCheckup.value.uuid, payload);

        if (newFiles.value.length > 0) {
            await uploadCheckupImages(selectedCheckup.value.uuid, newFiles.value);
        }

        emit('show-snackbar', 'Checkup updated successfully.', 'success');
        isEditDialogOpen.value = false;
        await loadCheckups();
    } catch (error) {
        emit('show-snackbar', 'Failed to update checkup.', 'error');
    }
}

async function deleteImage(image: ImageDto) {
    if (!selectedCheckup.value) return;
    const isConfirmed = await confirmDialog.value.Open({
        Title: 'Delete Image',
        Message: `Are you sure you want to delete this image? This action cannot be undone.`,
    });

    if (isConfirmed) {
        try {
            await deleteCheckupImage(selectedCheckup.value.uuid, image.uuid);
            existingImages.value = existingImages.value.filter(img => img.uuid !== image.uuid);
            emit('show-snackbar', 'Image deleted successfully.', 'success');
            await loadCheckups();
        } catch (error) {
            emit('show-snackbar', 'Failed to delete image.', 'error');
        }
    }
}

async function confirmAndDelete(checkup: CheckupDto) {
    const isConfirmed = await confirmDialog.value.Open({
        Title: 'Delete Checkup',
        Message: `Are you sure you want to delete the checkup from ${new Date(checkup.checkupDate).toLocaleDateString()}?`,
    });

    if (isConfirmed) {
        try {
            await deleteCheckup(checkup.uuid);
            emit('show-snackbar', 'Checkup deleted successfully.', 'success');
            await loadCheckups();
        } catch (error) {
            emit('show-snackbar', 'Failed to delete checkup.', 'error');
        }
    }
}

function openGallery(images: ImageDto[]) {
    selectedCheckupImages.value = images;
    isGalleryOpen.value = true;
}

watch(() => props.patient, (newPatient) => {
    if (newPatient) {
        loadCheckups();
    }
}, { immediate: true });
</script>