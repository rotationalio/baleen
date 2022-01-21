import React from 'react';
import { Navigate, useRoutes } from 'react-router-dom';
import Dashboard from 'pages/dashboard';
import { withBrowserRouter } from 'hoc';
import VerticalLayout from 'layouts/Vertical';

const Vocabulary = React.lazy(() => import('../pages/vocabulary'));
const Topics = React.lazy(() => import('../pages/topics'));

const AppRoutes = () => {
    const getLayout = () => {
        return VerticalLayout;
    };

    const Layout = getLayout();

    const Elements = useRoutes([
        {
            path: '/',
            element: <Layout />,
            children: [
                {
                    path: '',
                    element: <Navigate to={'dashboard'} />,
                },
                {
                    path: 'dashboard',
                    element: <Dashboard />,
                },
                {
                    path: 'vocabulary',
                    element: <Vocabulary />,
                },
                {
                    path: 'topics',
                    element: <Topics />,
                },
                {
                    path: '*',
                    element: <Navigate to={'/dashboard'} />,
                },
            ],
        },
    ]);

    return Elements;
};

export default withBrowserRouter(AppRoutes);
