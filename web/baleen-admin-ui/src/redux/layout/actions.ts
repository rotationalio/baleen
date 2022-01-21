// constants
import { LayoutActionTypes } from './constants';

export interface LayoutActionType<TPayload> {
    type:
        | LayoutActionTypes.CHANGE_LAYOUT
        | LayoutActionTypes.CHANGE_LAYOUT_WIDTH
        | LayoutActionTypes.CHANGE_MENU_POSITIONS
        | LayoutActionTypes.CHANGE_SIDEBAR_THEME
        | LayoutActionTypes.CHANGE_SIDEBAR_TYPE
        | LayoutActionTypes.TOGGLE_SIDEBAR_USER_INFO
        | LayoutActionTypes.CHANGE_TOPBAR_THEME
        | LayoutActionTypes.SHOW_RIGHT_SIDEBAR
        | LayoutActionTypes.HIDE_RIGHT_SIDEBAR;
    payload?: TPayload;
}

export const changeLayout = (layout: string): LayoutActionType<string> => ({
    type: LayoutActionTypes.CHANGE_LAYOUT,
    payload: layout,
});

export const changeLayoutWidth = (width: string): LayoutActionType<string> => ({
    type: LayoutActionTypes.CHANGE_LAYOUT_WIDTH,
    payload: width,
});

export const changeMenuPositions = (position: string): LayoutActionType<string> => ({
    type: LayoutActionTypes.CHANGE_MENU_POSITIONS,
    payload: position,
});

export const changeSidebarTheme = (sidebarTheme: string): LayoutActionType<string> => ({
    type: LayoutActionTypes.CHANGE_SIDEBAR_THEME,
    payload: sidebarTheme,
});

export const changeSidebarType = (sidebarType: string): LayoutActionType<string> => ({
    type: LayoutActionTypes.CHANGE_SIDEBAR_TYPE,
    payload: sidebarType,
});

export const toggleSidebarUserInfo = (showSidebarUserInfo: boolean): LayoutActionType<boolean> => ({
    type: LayoutActionTypes.TOGGLE_SIDEBAR_USER_INFO,
    payload: showSidebarUserInfo,
});

export const changeTopbarTheme = (topbarTheme: string): LayoutActionType<string> => ({
    type: LayoutActionTypes.CHANGE_TOPBAR_THEME,
    payload: topbarTheme,
});

export const showRightSidebar = (): LayoutActionType<null> => ({
    type: LayoutActionTypes.SHOW_RIGHT_SIDEBAR,
});

export const hideRightSidebar = (): LayoutActionType<null> => ({
    type: LayoutActionTypes.HIDE_RIGHT_SIDEBAR,
});
