import axios, { AxiosRequestConfig } from 'axios';
import Config from 'config';
import axiosErrorInterceptor from './error.interceptor';
import axiosResponseInterceptor from './response.interceptor';

const reqConfig: AxiosRequestConfig = {
    baseURL: Config.API_ENDPOINT,
    headers: {
        'Content-Type': 'application/json',
    },
};

const instance = axios.create(reqConfig);

axios.interceptors.response.use(axiosResponseInterceptor, axiosErrorInterceptor);

class apiCore {
    /**
     * Fetches data from given url
     */
    get = (url: string, params?: Record<string, string>, config?: AxiosRequestConfig) => {
        let response;
        if (params) {
            var queryString = params
                ? Object.keys(params)
                      .map((key) => key + '=' + params[key])
                      .join('&')
                : '';
            response = instance.get(`${url}?${queryString}`, config);
        } else {
            response = instance.get(`${url}`, config);
        }
        return response;
    };
}

export default apiCore;
