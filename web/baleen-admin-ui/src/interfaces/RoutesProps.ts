import { RouteProps } from 'react-router-dom';

export interface RoutesProps {
    path: RouteProps['path'];
    name?: string;
    element?: RouteProps['element'];
    route?: any;
    icon?: string;
    header?: string;
    roles?: string[];
    children?: RoutesProps[];
    index?: RouteProps['index'];
}
