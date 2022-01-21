import React from 'react';
import { BrowserRouter as Router } from 'react-router-dom';

const withBrowserRouter = (
    Component: React.FC<React.ReactElement<any, string | React.JSXElementConstructor<any>> | null>
) => {
    const Wrapper = (props: any) => {
        return (
            <Router>
                <Component {...props} />
            </Router>
        );
    };
    return Wrapper;
};

export { withBrowserRouter };
