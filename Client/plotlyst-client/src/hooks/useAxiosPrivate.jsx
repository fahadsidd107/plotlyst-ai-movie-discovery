import { useEffect, useRef } from 'react';
import axios from 'axios';
import useAuth from './useAuth';

const apiUrl = import.meta.env.VITE_API_BASE_URL;

// One stable client prevents consumers from re-running effects on every render.
const axiosAuth = axios.create({
    baseURL: apiUrl,
    withCredentials: true,
});

const useAxiosPrivate = () => {
    const { setAuth } = useAuth();
    const isRefreshing = useRef(false);
    const failedQueue = useRef([]);

    useEffect(() => {
        const processQueue = (error) => {
            failedQueue.current.forEach(({ resolve, reject }) => {
                if (error) reject(error);
                else resolve();
            });
            failedQueue.current = [];
        };

        const interceptor = axiosAuth.interceptors.response.use(
            response => response,
            async error => {
                const originalRequest = error.config;
                const status = error.response?.status;

                if (!originalRequest) return Promise.reject(error);
                if (originalRequest.url?.includes('/refresh') && status === 401) {
                    return Promise.reject(error);
                }

                if (status === 401 && !originalRequest._retry) {
                    if (isRefreshing.current) {
                        return new Promise((resolve, reject) => {
                            failedQueue.current.push({ resolve, reject });
                        }).then(() => axiosAuth(originalRequest));
                    }

                    originalRequest._retry = true;
                    isRefreshing.current = true;

                    try {
                        await axiosAuth.post('/refresh');
                        processQueue(null);
                        return axiosAuth(originalRequest);
                    } catch (refreshError) {
                        processQueue(refreshError);
                        localStorage.removeItem('user');
                        setAuth(null);
                        return Promise.reject(refreshError);
                    } finally {
                        isRefreshing.current = false;
                    }
                }

                return Promise.reject(error);
            }
        );

        return () => axiosAuth.interceptors.response.eject(interceptor);
    }, [setAuth]);

    return axiosAuth;
};

export default useAxiosPrivate;
