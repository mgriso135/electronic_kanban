import axios from 'axios';

// Changed from http://localhost:8080/api to relative routing for production
const API_BASE_URL = '/api'; 

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

// Function to handle successful responses
const handleSuccess = (response) => {
    return response;
};

// Function to handle error responses
const handleError = (error) => {
    console.error("API Error:", error);
    return Promise.reject(error);
};

// Add response interceptors
api.interceptors.response.use(handleSuccess, handleError);

export default api;