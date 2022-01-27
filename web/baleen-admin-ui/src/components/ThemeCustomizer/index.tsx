import React, { useState, useEffect, useCallback } from 'react';
import { useSelector, useDispatch } from 'react-redux';

// actions
import {
    changeLayout,
    changeLayoutWidth,
    changeMenuPositions,
    changeSidebarTheme,
    changeSidebarType,
    toggleSidebarUserInfo,
    changeTopbarTheme,
} from '../../redux/actions';

// store
import { RootState, AppDispatch } from '../../redux/store';

// constants
import * as layoutConstants from '../../constants/layout';

// custom hook
import { useViewport } from '../../hooks/useViewPort';

// components
import LayoutTypes from './LayoutTypes';
import LayoutWidth from './LayoutWidth';
import MenuPositions from './MenuPositions';
import LeftSideBarTheme from './LeftSideBarTheme';
import LeftSideBarType from './LeftSideBarType';
import SidebarUserInfo from './SidebarUserInfo';
import TopbarTheme from './TopbarTheme';

const ThemeCustomizer = () => {
    const dispatch = useDispatch<AppDispatch>();
    const { width } = useViewport();

    const {
        layoutType,
        layoutWidth,
        menuPosition,
        leftSideBarType,
        leftSideBarTheme,
        showSidebarUserInfo,
        topbarTheme,
    } = useSelector((state: RootState) => ({
        layoutType: state.Layout.layoutType,
        layoutWidth: state.Layout.layoutWidth,
        menuPosition: state.Layout.menuPosition,
        leftSideBarTheme: state.Layout.leftSideBarTheme,
        leftSideBarType: state.Layout.leftSideBarType,
        showSidebarUserInfo: state.Layout.showSidebarUserInfo,
        topbarTheme: state.Layout.topbarTheme,
    }));

    const [disableMenuPositions, setDisableMenuPositions] = useState<boolean>(false);
    const [disableSidebarTheme, setDisableSidebarTheme] = useState<boolean>(false);
    const [disableSidebarType, setDisableSidebarType] = useState<boolean>(false);
    const [disableSidebarUser, setDisableSidebarUser] = useState<boolean>(false);

    /**
     * change state based on props changes
     */
    const _loadStateFromProps = useCallback(() => {
        setDisableMenuPositions(layoutType !== layoutConstants.LayoutTypes.LAYOUT_TWO_COLUMN);
        setDisableSidebarTheme(layoutType !== layoutConstants.LayoutTypes.LAYOUT_HORIZONTAL);
        setDisableSidebarType(
            layoutType !== layoutConstants.LayoutTypes.LAYOUT_HORIZONTAL &&
                layoutType !== layoutConstants.LayoutTypes.LAYOUT_TWO_COLUMN &&
                width > 991
        );
        setDisableSidebarUser(
            layoutType !== layoutConstants.LayoutTypes.LAYOUT_HORIZONTAL &&
                layoutType !== layoutConstants.LayoutTypes.LAYOUT_TWO_COLUMN
        );
    }, [layoutType, width]);

    useEffect(() => {
        _loadStateFromProps();
    }, [_loadStateFromProps]);

    /**
     * On layout change
     */
    const changeLayoutType = (value: any) => {
        var layout = value;
        switch (layout) {
            case 'horizontal':
                dispatch(changeLayout(layoutConstants.LayoutTypes.LAYOUT_HORIZONTAL));
                break;
            case 'detached':
                dispatch(changeLayout(layoutConstants.LayoutTypes.LAYOUT_DETACHED));
                break;
            case 'vertical':
                dispatch(changeLayout(layoutConstants.LayoutTypes.LAYOUT_VERTICAL));
                break;
            default:
                dispatch(changeLayout(layoutConstants.LayoutTypes.LAYOUT_TWO_COLUMN));
                break;
        }
    };

    /**
     * Change the width mode
     */
    const changeWidthMode = (value: any) => {
        var mode = value;

        switch (mode) {
            case 'boxed':
                dispatch(changeLayoutWidth(layoutConstants.LayoutWidth.LAYOUT_WIDTH_BOXED));
                break;
            default:
                dispatch(changeLayoutWidth(layoutConstants.LayoutWidth.LAYOUT_WIDTH_FLUID));
                break;
        }
    };

    /*
     * Change menus position
     */
    const changeMenuPosition = (value: any) => {
        var position = value;

        switch (position) {
            case 'scrollable':
                dispatch(changeMenuPositions(layoutConstants.MenuPositions.MENU_POSITION_SCROLLABLE));
                break;
            default:
                dispatch(changeMenuPositions(layoutConstants.MenuPositions.MENU_POSITION_FIXED));
                break;
        }
    };

    /**
     * Changes the theme
     */
    const changeLeftSidebarTheme = (value: any) => {
        var theme = value;
        switch (theme) {
            case 'dark':
                dispatch(changeSidebarTheme(layoutConstants.SideBarTheme.LEFT_SIDEBAR_THEME_DARK));
                break;
            default:
                dispatch(changeSidebarTheme(layoutConstants.SideBarTheme.LEFT_SIDEBAR_THEME_LIGHT));
                break;
        }
    };

    /**
     * Change the leftsiderbar type
     */
    const changeLeftSiderbarType = (value: any) => {
        var type = value;
        switch (type) {
            case 'condensed':
                dispatch(changeSidebarType(layoutConstants.SideBarTypes.LEFT_SIDEBAR_TYPE_CONDENSED));
                break;
            case 'compact':
                dispatch(changeSidebarType(layoutConstants.SideBarTypes.LEFT_SIDEBAR_TYPE_COMPACT));
                break;
            default:
                dispatch(changeSidebarType(layoutConstants.SideBarTypes.LEFT_SIDEBAR_TYPE_DEFAULT));
                break;
        }
    };

    /*
     * Toggle the leftsidebar use info
     */
    const toggleLeftSidebarUserInfo = (value: any) => {
        var checked = value;
        dispatch(toggleSidebarUserInfo(checked));
    };

    /*
     * Change topbar theme
     */
    const changeTopBarTheme = (value: any) => {
        var theme = value;
        switch (theme) {
            case 'light':
                dispatch(changeTopbarTheme(layoutConstants.TopbarTheme.TOPBAR_THEME_LIGHT));
                break;
            default:
                dispatch(changeTopbarTheme(layoutConstants.TopbarTheme.TOPBAR_THEME_DARK));
                break;
        }
    };

    /**
     * Reset everything
     */
    const reset = () => {
        changeLayoutType(layoutConstants.LayoutTypes.LAYOUT_TWO_COLUMN);
        changeWidthMode(layoutConstants.LayoutWidth.LAYOUT_WIDTH_FLUID);
        changeMenuPosition(layoutConstants.MenuPositions.MENU_POSITION_FIXED);
        changeLeftSidebarTheme(layoutConstants.SideBarTheme.LEFT_SIDEBAR_THEME_LIGHT);
        toggleLeftSidebarUserInfo(false);
        changeTopBarTheme(layoutConstants.TopbarTheme.TOPBAR_THEME_LIGHT);
        changeLeftSiderbarType(layoutConstants.SideBarTypes.LEFT_SIDEBAR_TYPE_DEFAULT);
    };

    return (
        <React.Fragment>
            <h6 className="fw-medium px-3 m-0 py-2 text-uppercase bg-light">
                <span className="d-block py-1">Theme Settings</span>
            </h6>
            <div className="p-3">
                <div className="alert alert-warning" role="alert">
                    <strong>Customize </strong> the overall color scheme, sidebar menu, etc.
                </div>

                {/* Layouts */}
                <LayoutTypes
                    changeLayoutType={changeLayoutType}
                    layoutType={layoutType}
                    layoutConstants={layoutConstants.LayoutTypes}
                />

                {/* Width */}
                <LayoutWidth
                    changeWidthMode={changeWidthMode}
                    layoutWidth={layoutWidth}
                    layoutConstants={layoutConstants.LayoutWidth}
                />

                {/* Menu Posiotions */}
                {disableMenuPositions && (
                    <MenuPositions
                        menuPosition={menuPosition}
                        changeMenuPosition={changeMenuPosition}
                        layoutConstants={layoutConstants.MenuPositions}
                    />
                )}

                {/* Left Sidebar */}
                {disableSidebarTheme && (
                    <LeftSideBarTheme
                        changeLeftSidebarTheme={changeLeftSidebarTheme}
                        leftSideBarTheme={leftSideBarTheme}
                        layoutConstants={layoutConstants.SideBarTheme}
                    />
                )}

                {/* Left Sidebar Size */}
                {disableSidebarType && (
                    <LeftSideBarType
                        changeLeftSiderbarType={changeLeftSiderbarType}
                        leftSideBarType={leftSideBarType}
                        layoutConstants={layoutConstants.SideBarTypes}
                    />
                )}

                {/* User Info */}
                {disableSidebarUser && (
                    <SidebarUserInfo
                        toggleLeftSidebarUserInfo={toggleLeftSidebarUserInfo}
                        showSidebarUserInfo={showSidebarUserInfo}
                    />
                )}

                {/* Topbar */}
                <TopbarTheme
                    changeTopBarTheme={changeTopBarTheme}
                    topbarTheme={topbarTheme}
                    layoutConstants={layoutConstants.TopbarTheme}
                />

                <div className="d-grid mt-4">
                    <button className="btn btn-primary" id="resetBtn" onClick={() => reset()}>
                        Reset to Default
                    </button>
                </div>
            </div>
        </React.Fragment>
    );
};

export default ThemeCustomizer;
