// action constants
import { LayoutActionTypes, LayoutStateTypes } from './constants';

// app constants
import {
    LayoutTypes,
    LayoutWidth,
    MenuPositions,
    SideBarTheme,
    SideBarTypes,
    TopbarTheme,
} from '../../constants/layout';

// actions
import { LayoutActionType } from './actions';

// utils
import { getLayoutConfigs } from '../../utils';

const INIT_STATE = {
    layoutType: LayoutTypes.LAYOUT_VERTICAL,
    layoutWidth: LayoutWidth.LAYOUT_WIDTH_FLUID,
    menuPosition: MenuPositions.MENU_POSITION_FIXED,
    leftSideBarTheme: SideBarTheme.LEFT_SIDEBAR_THEME_LIGHT,
    leftSideBarType: SideBarTypes.LEFT_SIDEBAR_TYPE_DEFAULT,
    showSidebarUserInfo: false,
    topbarTheme: TopbarTheme.TOPBAR_THEME_LIGHT,
    isOpenRightSideBar: false,
};

const Layout = (state: LayoutStateTypes = INIT_STATE, action: LayoutActionType<string | boolean | null>) => {
    switch (action.type) {
        case LayoutActionTypes.CHANGE_LAYOUT:
            return {
                ...state,
                layoutType: action.payload,
            };
        case LayoutActionTypes.CHANGE_LAYOUT_WIDTH:
            const layoutConfig = getLayoutConfigs(action.payload! && action.payload);
            return {
                ...state,
                layoutWidth: action.payload,
                ...layoutConfig,
            };
        case LayoutActionTypes.CHANGE_MENU_POSITIONS:
            return {
                ...state,
                menuPosition: action.payload,
            };
        case LayoutActionTypes.CHANGE_SIDEBAR_THEME:
            return {
                ...state,
                leftSideBarTheme: action.payload,
            };
        case LayoutActionTypes.CHANGE_SIDEBAR_TYPE:
            return {
                ...state,
                leftSideBarType: action.payload,
            };
        case LayoutActionTypes.TOGGLE_SIDEBAR_USER_INFO:
            return {
                ...state,
                showSidebarUserInfo: action.payload,
            };
        case LayoutActionTypes.CHANGE_TOPBAR_THEME:
            return {
                ...state,
                topbarTheme: action.payload,
            };
        case LayoutActionTypes.SHOW_RIGHT_SIDEBAR:
            return {
                ...state,
                isOpenRightSideBar: true,
            };
        case LayoutActionTypes.HIDE_RIGHT_SIDEBAR:
            return {
                ...state,
                isOpenRightSideBar: false,
            };
        default:
            return state;
    }
};

export default Layout;
