import { defineStore } from 'pinia'

const stateFunction = () => {
    return {
        accessToken: null as string | null,
    };
};

const mySettingsBox = { state: stateFunction }

export const useAuthStore = defineStore('auth', mySettingsBox);