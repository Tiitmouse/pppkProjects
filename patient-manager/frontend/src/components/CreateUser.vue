<template>
  <v-container>
    <v-card max-width="800" class="mx-auto">
      <v-card-title class="text-h5">Create new</v-card-title>
      <v-card-text>
        <v-form ref="form" v-model="valid">
          <v-radio-group 
            v-if="user?.role === UserRole.SuperAdmin" 
            v-model="roleSelection" 
            inline 
            label="User Type"
          >
            <v-radio label="Doctor" value="doctor"></v-radio>
            <v-radio label="Patient" value="patient"></v-radio>
          </v-radio-group>

          <div v-if="roleSelection === 'doctor'">
            <v-text-field v-model="doctorForm.firstName" label="First Name" :rules="[rules.required]"></v-text-field>
            <v-text-field v-model="doctorForm.lastName" label="Last Name" :rules="[rules.required]"></v-text-field>
            <v-text-field v-model="doctorForm.oib" label="OIB" :rules="[rules.required]"></v-text-field>
            <v-text-field v-model="doctorForm.email" label="E-mail"
              :rules="[rules.required, rules.email]"></v-text-field>
          </div>

          <div v-if="roleSelection === 'patient'">
            <v-text-field v-model="patientForm.firstName" label="First Name" :rules="[rules.required]"></v-text-field>
            <v-text-field v-model="patientForm.lastName" label="Last Name" :rules="[rules.required]"></v-text-field>
            <v-text-field v-model="patientForm.oib" label="OIB" :rules="[rules.required]"></v-text-field>
            <v-text-field v-model="patientForm.birthDate" label="Birth Date" type="date"
              :rules="[rules.required]"></v-text-field>
            <v-radio-group v-model="patientForm.gender" label="Gender" :rules="[rules.required]" inline>
              <v-radio label="Male" value="m"></v-radio>
              <v-radio label="Female" value="f"></v-radio>
            </v-radio-group>
          </div>
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue-darken-1" variant="text" @click="cancel">Cancel</v-btn>
        <v-btn color="blue-darken-1" :disabled="!valid" @click="submit">Create</v-btn>
      </v-card-actions>
    </v-card>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import * as userService from '@/services/userService';
import * as patientService from '@/services/patientService';
import type { User } from '@/models/user';
import type { NewPatientDto } from '@/dtos/patientDto';
import { UserRole } from '@/enums/userRole';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const user = ref<User>();

const router = useRouter();
const valid = ref(false);
const roleSelection = ref('');
const today = new Date().toISOString().substr(0, 10);

const doctorForm = ref<User>({
  uuid: '',
  firstName: '',
  lastName: '',
  oib: '',
  birthDate: today,
  email: '',
  gender: '',
  role: UserRole.Doctor,
});

const password = ref('');

const patientForm = ref<NewPatientDto>({
  firstName: '',
  lastName: '',
  oib: '',
  birthDate: today,
  gender: '',
  uuid: '',
  role: UserRole.Patient
});

const rules = {
  required: (v: string) => !!v || 'This field is required.',
  email: (v: string) => /.+@.+\..+/.test(v) || 'E-mail must be valid.',
};

const cancel = () => {
  router.push('/users');
};

const submit = async () => {
  try {
    if (roleSelection.value === 'doctor') {
      await userService.createUser(doctorForm.value, password.value);
    } else {
      await patientService.createPatient(patientForm.value);
    }
    router.push('/users');
  } catch (error) {
    console.error('Failed to create:', error);
  }
};


onMounted(() => {
  const currentUser = authStore.User;
  user.value = currentUser;

  if (currentUser?.role === UserRole.Doctor) {
    roleSelection.value = 'patient';
  } 
  else if (currentUser?.role === UserRole.SuperAdmin) {
    roleSelection.value = 'doctor';
  }
});
</script>