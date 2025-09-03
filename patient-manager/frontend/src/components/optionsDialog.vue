<template>
    <v-dialog v-model="dialog" :max-width="width" persistent>
        <v-card>
            <v-card-title class="text-h5">
                {{ title }}
            </v-card-title>

            <v-list>
                <v-list-item v-for="(option, index) in options" :key="index" @click="select(option)">
                    <v-list-item-title>{{ option }}</v-list-item-title>
                </v-list-item>
            </v-list>

            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="error" variant="elevated" @click="cancel">
                    Cancel
                </v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script lang="ts" setup>
import { ref } from "vue";

interface OptionsProps {
    Title: string;
    Options: string[];
    Width?: number;
}

const dialog = ref(false);
const title = ref("");
const options = ref<string[]>([]);
const width = ref(400);

const resolve = ref<(value: string | null) => void>();

function select(option: string) {
    if (resolve.value) {
        resolve.value(option);
    }
    dialog.value = false;
}

function cancel() {
    if (resolve.value) {
        resolve.value(null);
    }
    dialog.value = false;
}

function Open({ Title, Options, Width = 400 }: OptionsProps): Promise<string | null> {
    dialog.value = true;
    title.value = Title;
    options.value = Options;
    width.value = Width;

    return new Promise((res) => {
        resolve.value = res;
    });
}

defineExpose({
    Open,
});
</script>