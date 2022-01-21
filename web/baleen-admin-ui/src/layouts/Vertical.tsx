import React, { Suspense, useCallback, useEffect, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { Container } from 'react-bootstrap';

// actions
import { changeSidebarType } from '../redux/actions';

// store
import { RootState, AppDispatch } from '../redux/store';

// constants
import { LayoutTypes, SideBarTypes } from '../constants';

// utils
import { changeBodyAttribute } from '../utils';
import { Outlet } from 'react-router-dom';

// code splitting and lazy loading
// https://blog.logrocket.com/lazy-loading-components-in-react-16-6-6cea535c0b52
const Topbar = React.lazy(() => import('./Topbar'));
const LeftSidebar = React.lazy(() => import('./LeftSidebar'));
const Footer = React.lazy(() => import('./Footer'));
const RightSidebar = React.lazy(() => import('./RightSidebar'));

const loading = () => <div className=""></div>;

interface VerticalLayoutProps {
    children?: React.ReactNode;
}

const VerticalLayout = ({ children }: VerticalLayoutProps) => {
    const dispatch = useDispatch<AppDispatch>();

    const {
        layoutWidth,
        menuPosition,
        leftSideBarTheme,
        leftSideBarType,
        showSidebarUserInfo,
        showTwoToneIcons,
        topbarTheme,
        isOpenRightSideBar,
    } = useSelector((state: RootState) => ({
        layoutWidth: state.Layout.layoutWidth,
        menuPosition: state.Layout.menuPosition,
        leftSideBarTheme: state.Layout.leftSideBarTheme,
        leftSideBarType: state.Layout.leftSideBarType,
        showSidebarUserInfo: state.Layout.showSidebarUserInfo,
        showTwoToneIcons: state.Layout.showTwoToneIcons,
        topbarTheme: state.Layout.topbarTheme,
        isOpenRightSideBar: state.Layout.isOpenRightSideBar,
    }));

    const [isMenuOpened, setIsMenuOpened] = useState<boolean>(false);

    /*
    layout defaults
    */
    useEffect(() => {
        changeBodyAttribute('data-layout-mode', LayoutTypes.LAYOUT_VERTICAL);
    }, []);

    useEffect(() => {
        changeBodyAttribute('data-layout-width', layoutWidth);
    }, [dispatch, layoutWidth]);

    useEffect(() => {
        changeBodyAttribute('data-layout-menu-position', menuPosition);
    }, [menuPosition]);

    useEffect(() => {
        changeBodyAttribute('data-sidebar-color', leftSideBarTheme);
    }, [leftSideBarTheme]);

    useEffect(() => {
        changeBodyAttribute('data-sidebar-size', leftSideBarType);
    }, [leftSideBarType]);

    useEffect(() => {
        changeBodyAttribute('data-topbar-color', topbarTheme);
    }, [topbarTheme]);

    useEffect(() => {
        changeBodyAttribute('data-sidebar-icon', showTwoToneIcons ? 'twotones' : '');
    }, [showTwoToneIcons]);

    useEffect(() => {
        changeBodyAttribute('data-sidebar-showuser', showSidebarUserInfo);
    }, [showSidebarUserInfo]);

    /**
     * Open the menu when having mobile screen
     */
    const openMenu = () => {
        setIsMenuOpened((prevState) => !prevState);

        if (document.body) {
            if (isMenuOpened) {
                document.body.classList.remove('sidebar-enable');
            } else {
                document.body.classList.add('sidebar-enable');
            }
        }
    };

    const updateDimensions = useCallback(() => {
        // activate the condensed sidebar if smaller devices like ipad or tablet
        if (window.innerWidth > 768 && window.innerWidth <= 1028) {
            dispatch(changeSidebarType(SideBarTypes.LEFT_SIDEBAR_TYPE_CONDENSED));
        } else if (window.innerWidth > 1028) {
            dispatch(changeSidebarType(SideBarTypes.LEFT_SIDEBAR_TYPE_DEFAULT));
        }
    }, [dispatch]);

    useEffect(() => {
        window.addEventListener('resize', updateDimensions);

        return () => {
            window.removeEventListener('resize', updateDimensions);
        };
    }, [dispatch, updateDimensions]);

    const isCondensed: boolean = leftSideBarType === SideBarTypes.LEFT_SIDEBAR_TYPE_CONDENSED;

    return (
        <>
            <div id="wrapper">
                <Suspense fallback={loading()}>
                    <Topbar openLeftMenuCallBack={openMenu} hideLogo={false} />
                </Suspense>
                <Suspense fallback={loading()}>
                    <LeftSidebar isCondensed={isCondensed} />
                </Suspense>
                <div className="content-page">
                    <div className="content">
                        <Container fluid>
                            <Suspense fallback={loading()}>
                                <Outlet />
                            </Suspense>
                        </Container>
                    </div>

                    <Suspense fallback={loading()}>
                        <Footer />
                    </Suspense>
                </div>
            </div>
            {isOpenRightSideBar && (
                <Suspense fallback={loading()}>
                    <RightSidebar />
                </Suspense>
            )}
        </>
    );
};
export default VerticalLayout;
