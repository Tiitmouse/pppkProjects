<template>
  <div>
    <v-container class="d-flex align-center justify-center">
      <v-card width="800" rounded="lg" elevation="3" class="pa-8">
        <div>
          <div class="text-h5 mr-4">Create new</div>
          <v-radio-group v-model="selectedRole" inline hide-details class="mt-0 pt-0">
            <v-radio label="Doctor" value="Doctor"></v-radio>
            <v-radio label="Patient" value="Patient"></v-radio>
          </v-radio-group>
        </div>

        <!-- Common Fields for Both Roles -->
        <v-text-field v-model="user.firstName" :error-messages="errors.firstName" label="First Name" variant="outlined"
          density="compact" @input="updateErrors('firstName', '')"></v-text-field>

        <v-text-field v-model="user.lastName" :error-messages="errors.lastName" label="Last Name" variant="outlined"
          density="compact" class="mt-2" @input="updateErrors('lastName', '')"></v-text-field>

        <v-text-field v-model="user.email" :error-messages="errors.email" label="Email Address" type="email"
          variant="outlined" density="compact" class="mt-2" @input="updateErrors('email', '')"
          autocomplete="off"></v-text-field>

        <v-text-field v-model="password" :error-messages="errors.password" label="Password" type="password"
          variant="outlined" density="compact" class="mt-2" @input="updateErrors('password', '')"
          autocomplete="off"></v-text-field>

        <!-- Patient-Specific Fields -->
        <div v-if="selectedRole === 'Patient'">
          <v-text-field v-model="user.oib" :error-messages="errors.oib" label="OIB" variant="outlined" density="compact"
            class="mt-2" @input="updateErrors('oib', '')"></v-text-field>

          <v-text-field v-model="user.residence" :error-messages="errors.residence" label="Residence" variant="outlined"
            density="compact" class="mt-2" @input="updateErrors('residence', '')"></v-text-field>

          <v-text-field v-model="user.birthDate" :error-messages="errors.birthDate" label="Birth Date" type="date"
            variant="outlined" density="compact" class="mt-2" @input="updateErrors('birthDate', '')"></v-text-field>
        </div>

        <div class="d-flex gap-4 mt-6">
          <v-btn color="primary" size="large" block :loading="isSubmitting" :disabled="isSubmitting"
            @click="submitForm">
            {{ isSubmitting ? 'Creating...' : 'Create User' }}
          </v-btn>

        </div>

        <v-alert v-if="errorMessage" type="error" variant="tonal" class="mt-6" closable
          @click:close="errorMessage = ''">
          {{ errorMessage }}
        </v-alert>

        <v-alert v-if="successMessage" type="success" variant="tonal" class="mt-6" closable
          @click:close="successMessage = ''">
          {{ successMessage }}
        </v-alert>

      </v-card>
    </v-container>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, watch } from 'vue';
import type { User } from '@/models/user';
import { UserRole } from '@/enums/userRole';
import type { FormErrors } from '@/models/formErrors';
import { isOibValid } from '@/utils/validateOIB';
import { isEmailValid } from '@/utils/validateEmail';
import { createUser } from '@/services/userService';

const createInitialUser = (): User => ({
  uuid: '',
  firstName: '',
  lastName: '',
  oib: '',
  residence: '',
  birthDate: '',
  email: '',
  role: UserRole.Doctor
});

const selectedRole = ref<UserRole>(UserRole.Doctor);

const user = ref<User>(createInitialUser());
const password = ref('');

const isSubmitting = ref(false);
const errorMessage = ref('');
const successMessage = ref('');
const errors = reactive<FormErrors>({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  oib: '',
  residence: '',
  birthDate: '',
});

const updateErrors = (field: keyof FormErrors, value: string) => {
  if (field in errors) {
    errors[field] = value;
  }
};

const validateForm = (): boolean => {
  Object.keys(errors).forEach(key => { errors[key as keyof FormErrors] = '' });
  let isValid = true;

  // Common validations
  if (!user.value.firstName.trim()) { errors.firstName = 'First name is required'; isValid = false; }
  if (!user.value.lastName.trim()) { errors.lastName = 'Last name is required'; isValid = false; }
  if (!user.value.email.trim()) { errors.email = 'Email is required'; isValid = false; }
  else if (!isEmailValid(user.value.email)) { errors.email = 'Enter a valid email address'; isValid = false; }
  if (!password.value) { errors.password = 'Password is required'; isValid = false; }
  else if (password.value.length < 6) { errors.password = 'Password must be at least 6 characters'; isValid = false; }

  // Patient-specific validations
  if (selectedRole.value === UserRole.Patient) {
    if (!user.value.oib.trim()) { errors.oib = 'OIB is required'; isValid = false; }
    else if (!isOibValid(user.value.oib)) { errors.oib = 'Invalid OIB'; isValid = false; }
    if (!user.value.residence.trim()) { errors.residence = 'Residence is required'; isValid = false; }
    if (!user.value.birthDate) { errors.birthDate = 'Birth date is required'; isValid = false; }
  }

  return isValid;
};

const submitForm = async () => {
  errorMessage.value = '';
  successMessage.value = '';

  if (!validateForm()) {
    return;
  }

  isSubmitting.value = true;
  try {
    user.value.role = selectedRole.value;

    const createdUser = await createUser(user.value, password.value);

    if (createdUser) {
      successMessage.value = `User ${createdUser.email} created successfully!`;
      resetForm();
    }
  } catch (error: any) {
    console.error('Error creating user:', error);
    errorMessage.value = error.response?.data?.message || 'Failed to create user. Please try again.';
  } finally {
    isSubmitting.value = false;
  }
};

const resetForm = () => {
  user.value = createInitialUser();
  password.value = '';
  Object.keys(errors).forEach(key => { errors[key as keyof FormErrors] = '' });
};

watch(selectedRole, (newRole) => {
  resetForm();
  user.value.role = newRole;
});
</script>

<style lang="css" scoped>
.gap-4 {
  gap: 1rem;
}
</style>
