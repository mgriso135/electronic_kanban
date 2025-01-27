import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api'; // Adjust if your backend runs on a different port or URL

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