import { MENU_ITEMS, MenuItemTypes } from '../constants/menu';

const getMenuItems = () => {
    return MENU_ITEMS;
};

const findAllParent = (menuItems: MenuItemTypes[], menuItem: MenuItemTypes): string[] => {
    let parents: string[] = [];
    const parent = findMenuItem(menuItems, menuItem['parentKey']);

    if (parent) {
        parents.push(parent['key']);
        if (parent['parentKey']) parents = [...parents, ...findAllParent(menuItems, parent)];
    }

    return parents;
};

const findMenuItem = (
    menuItems: MenuItemTypes[] | undefined,
    menuItemKey: MenuItemTypes['key'] | undefined
): MenuItemTypes | null => {
    if (menuItems && menuItemKey) {
        for (var i = 0; i < menuItems.length; i++) {
            if (menuItems[i].key === menuItemKey) {
                return menuItems[i];
            }
            var found = findMenuItem(menuItems[i].children, menuItemKey);
            if (found) return found;
        }
    }
    return null;
};

export { getMenuItems, findAllParent, findMenuItem };
