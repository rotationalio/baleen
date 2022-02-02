import { AxiosResponse } from 'axios';

const axiosResponseInterceptor = (response: AxiosResponse) => {
    return response;
};

export default axiosResponseInterceptor;
