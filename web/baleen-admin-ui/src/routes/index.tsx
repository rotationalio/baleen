import { Navigate, useRoutes } from 'react-router-dom';
import VerticalLayout from 'layouts/Vertical';
import Topics from 'pages/topics';
import Vocabulary from 'pages/vocabulary';
import Dashboard from 'pages/dashboard';
import { withBrowserRouter } from 'hoc';

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
