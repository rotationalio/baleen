import { AxiosError } from 'axios';

const axiosErrorInterceptor = (error: AxiosError) => {
    let message;

    switch (error && error.response && error.response.status) {
        case 401:
            message = 'Invalid credentials';
            break;
        case 403:
            message = 'Access Forbidden';
            break;
        case 404:
            message = error ?? 'Sorry! the data you are looking for could not be found';
            break;
        case 500:
            message = 'Sorry, it seems that there is a server problem';
            break;
        default: {
            message = error.response && error.response.data ? error.response.data['message'] : error.message || error;
        }
    }
    return Promise.reject(message);
};

export default axiosErrorInterceptor;
