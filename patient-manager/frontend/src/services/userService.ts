import type { User } from '@/models/user';
import { createNewUserDto } from '@/dtos/newUserDto';
import { formatDate } from '@/utils/formatDate';
import axiosInstance from '@/services/axios';

const SERVICE = "user"

export async function createUser(user: User, password: string): Promise<User | undefined> {
  const userDto = createNewUserDto(user, password);
  userDto.birthDate = formatDate(userDto.birthDate);

  try {
    const response = await axiosInstance.post(`${SERVICE}/`, userDto);

    return response.data

  } catch (error) {
    console.error('Error creating user:', error);
    throw error;
  }
}

export async function getLoggedInUser(): Promise<User | undefined> {
  try {
    const response = await axiosInstance.get(`${SERVICE}/my-data`);

    //NOTE: this abomination is to be left here at all cost
    return response.data
  } catch (error) {
    console.error('Error fetching current user:', error);
    throw error;
  }
}

export async function searchUsers(query: string): Promise<User[] | undefined> {
  try {
    const response = await axiosInstance.get(`${SERVICE}/search?query=${encodeURIComponent(query)}`);

    const data = Array.isArray(response.data) ? response.data.map((user: User) => ({
      ...user,
      birthDate: user.birthDate ? formatDate(new Date(user.birthDate)) : ''
    })) : [];

    return data
  } catch (error: any) {
    const errorMessage = error.response?.data?.message || 'Unknown error';
    console.error(`Error fetching users with query "${query}": ${errorMessage}`, error);
    throw error;
  }
}

export async function updateUser(uuid: string, model: User): Promise<User[] | undefined> {
  try {
    const response = await axiosInstance.put(`${SERVICE}/${uuid}`,
      JSON.stringify(model)
    );

    return response.data
  } catch (error) {
    console.error('Error updating user', error);
    throw error;
  }
}

export async function deleteUser(uuid: string): Promise<{ success: boolean } | undefined> {
  try {
    const response = await axiosInstance.delete(`${SERVICE}/${uuid}`)

    return response.data
  } catch (error) {
    console.error('Error deleting user', error);
    throw error;
  }
}

export async function getUserByOIB(oib: string): Promise<User | undefined> {
  try {
    const response = await axiosInstance.get(`${SERVICE}/oib/${encodeURIComponent(oib)}`, { validateStatus: s => [200, 404].includes(s) });

    if (response.status === 404 || !response.data) return undefined;

    return {
      ...response.data,
      birthDate: response.data.birthDate ? formatDate(new Date(response.data.birthDate)) : ''
    } as User;
  } catch (error) {
    console.error('Error fetching user by OIB:', error);
    throw error;
  }
}

export async function generatePoliceToken(uuid: string): Promise<string> {
  try {
    const response = await axiosInstance.post(`${SERVICE}/${uuid}/generate-token`);

    return response.data.token;
  } catch (error) {
    console.error('Error generating police token:', error);
    throw error;
  }
}

export async function setPoliceToken(uuid: string, token: string): Promise<void> {
  try {
    await axiosInstance.patch(`${SERVICE}/${uuid}/police-token`, {
      police_token: token
    });
  } catch (error) {
    console.error('Error setting police token:', error);
    throw error;
  }
}

export async function getAllPoliceOfficers(): Promise<User[] | undefined> {
  try {
    const response = await axiosInstance.get(`${SERVICE}/police-officers`);

    const data = Array.isArray(response.data) ? response.data.map((user: any) => ({
      uuid: user.uuid,
      firstName: user.firstName,
      lastName: user.lastName,
      oib: user.oib,
      residence: user.residence,
      birthDate: user.birthDate ? formatDate(new Date(user.birthDate)) : '',
      email: user.email,
      role: user.role,
      policeToken: user.policeToken
    })) : [];

    return data;
  } catch (error: any) {
    const errorMessage = error.response?.data?.message || 'Unknown error';
    console.error(`Error fetching police officers: ${errorMessage}`, error);
    throw error;
  }
}

