import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import classNames from 'classnames';
import FeatherIcon from 'feather-icons-react';

// actions
import { changeSidebarType } from 'redux/actions';

// store
import { RootState, AppDispatch } from 'redux/store';

//constants
import { LayoutTypes, SideBarTypes } from 'constants/layout';

export interface NotificationItem {
    id: number;
    text: string;
    subText: string;
    icon?: string;
    avatar?: string;
    bgColor?: string;
}

interface TopbarProps {
    hideLogo?: boolean;
    navCssClasses?: string;
    openLeftMenuCallBack?: () => void;
    topbarDark?: boolean;
}

const Topbar = ({ hideLogo, navCssClasses, openLeftMenuCallBack }: TopbarProps) => {
    const dispatch = useDispatch<AppDispatch>();

    const [isopen, setIsopen] = useState<boolean>(false);

    const navbarCssClasses: string = navCssClasses || '';
    const containerCssClasses: string = !hideLogo ? 'container-fluid' : '';

    const { layoutType, leftSideBarType } = useSelector((state: RootState) => ({
        layoutType: state.Layout.layoutType,
        leftSideBarType: state.Layout.leftSideBarType,
    }));

    /**
     * Toggle the leftmenu when having mobile screen
     */
    const handleLeftMenuCallBack = () => {
        setIsopen(!isopen);
        if (openLeftMenuCallBack) openLeftMenuCallBack();
    };

    /**
     * Toggles the left sidebar width
     */
    const toggleLeftSidebarWidth = () => {
        if (leftSideBarType === 'default' || leftSideBarType === 'compact')
            dispatch(changeSidebarType(SideBarTypes.LEFT_SIDEBAR_TYPE_CONDENSED));
        if (leftSideBarType === 'condensed') dispatch(changeSidebarType(SideBarTypes.LEFT_SIDEBAR_TYPE_DEFAULT));
    };

    return (
        <React.Fragment>
            <div className={`navbar-custom ${navbarCssClasses}`}>
                <div className={containerCssClasses}>
                    {!hideLogo && <div className="logo-box"></div>}

                    <ul className="list-unstyled topnav-menu topnav-menu-left m-0">
                        {layoutType !== LayoutTypes.LAYOUT_HORIZONTAL && (
                            <li>
                                <button
                                    className="button-menu-mobile d-none d-lg-block"
                                    onClick={toggleLeftSidebarWidth}>
                                    <FeatherIcon icon="menu" />
                                    <i className="fe-menu"></i>
                                </button>
                            </li>
                        )}

                        <li>
                            <button className="button-menu-mobile d-lg-none d-bolck" onClick={handleLeftMenuCallBack}>
                                <FeatherIcon icon="menu" />
                            </button>
                        </li>

                        {/* Mobile menu toggle (Horizontal Layout) */}
                        <li>
                            <Link
                                to="#"
                                className={classNames('navbar-toggle nav-link', {
                                    open: isopen,
                                })}
                                onClick={handleLeftMenuCallBack}>
                                <div className="lines">
                                    <span></span>
                                    <span></span>
                                    <span></span>
                                </div>
                            </Link>
                        </li>
                    </ul>
                </div>
            </div>
        </React.Fragment>
    );
};

export default Topbar;
